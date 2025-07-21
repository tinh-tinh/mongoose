package mongoose

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseSchema struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type BaseTimestamp struct {
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type ModelCommon interface {
	SetConnect(connect *Connect)
	GetName() string
}

type Model[M any] struct {
	option     *ModelOptions
	docs       []bson.E
	connect    *Connect
	idx        mongo.IndexModel
	preHooks   []Hook[M]
	postHooks  []Hook[M]
	Ctx        context.Context
	Collection *mongo.Collection
}

type ModelOptions struct {
	Timestamp  bool
	ID         bool
	Validation bool
}

// NewModel returns a new instance of Model[M] with the given connect and name
// name is the name of the collection in the database
// the returned Model[M] is used to interact with the collection in the database
func NewModel[M any](opts ...ModelOptions) *Model[M] {
	defaultOption := ModelOptions{
		ID:         true,
		Timestamp:  true,
		Validation: true,
	}

	if len(opts) > 0 {
		defaultOption = common.MergeStruct(opts...)
	}

	return &Model[M]{
		option: &defaultOption,
	}
}

// SetConnect sets the context and collection of the model to the given connect.
// It is used internally by the ForFeature function to set the connect of the model.
// The given connect must be a *Connect.
func (m *Model[M]) SetConnect(connect *Connect) {
	m.Ctx = connect.Ctx
	m.connect = connect
	m.Collection = connect.Client.Database(connect.DB).Collection(m.GetName())
	if !reflect.ValueOf(m.idx).IsZero() {
		_, err := m.Collection.Indexes().CreateOne(m.Ctx, m.idx)
		if err != nil {
			log.Println(err)
		}
	}
}

func (m *Model[M]) SetContext(ctx context.Context) {
	m.Ctx = ctx
}

// GetName returns the name of the collection in the database
func (m *Model[M]) GetName() string {
	var model M
	ctModel := reflect.ValueOf(&model).Elem()

	fnc := ctModel.MethodByName("CollectionName")
	var name string
	if fnc.IsValid() {
		name = fnc.Call(nil)[0].String()
	} else {
		name = common.GetStructName(model)
	}
	return name
}

func (m *Model[M]) Index(idx bson.D, unique bool) {
	indexModel := mongo.IndexModel{
		Keys:    idx,
		Options: options.Index().SetUnique(unique),
	}
	m.idx = indexModel
}

// Set sets the data of the model to the given data.
// It iterates over the struct fields of the given data and sets the corresponding field of the model to the value of the field.
// If the field is not found in the model, it is not set.
// If the field is not tagged with bson, the field name is used as the key.
// If the field is tagged with bson, the tag value is used as the key.
// The given data must be a struct.
func (m *Model[M]) Set(data interface{}) {
	var model M

	ctInput := reflect.ValueOf(data).Elem()
	for i := range ctInput.NumField() {
		name := ctInput.Type().Field(i).Name
		val := ctInput.Field(i).Interface()

		if val != nil {
			ctModel := reflect.ValueOf(&model).Elem()

			if ctModel.FieldByName(name).IsValid() {
				field, ok := ctModel.Type().FieldByName(name)
				if ok {
					key := field.Tag.Get("bson")
					m.docs = append(m.docs, bson.E{Key: key, Value: val})
				}
			}
		}
	}
}

// Save saves the changes to the model in the database.
// If the model has no ID, InsertOne is used to insert the document.
// If the model has an ID, UpdateByID is used to update the existing document.
// The createdAt and updatedAt fields are automatically set if not present.
// If the model has no changes, Save does nothing.
// Save returns an error if the operation fails.
func (m *Model[M]) Save() error {
	err := ExecutePreHook(Save, m, m.docs)
	if err != nil {
		return err
	}

	if len(m.docs) == 0 {
		return nil
	}

	idIndex := slices.IndexFunc(m.docs, func(e bson.E) bool {
		return e.Key == "_id"
	})
	if idIndex == -1 {
		inserts := m.docs

		if m.option.ID {
			inserts = append(m.docs,
				bson.E{Key: "_id", Value: primitive.NewObjectID()},
			)
		}

		if m.option.Timestamp {
			inserts = append(m.docs,
				bson.E{Key: "createdAt", Value: time.Now()},
				bson.E{Key: "updatedAt", Value: time.Now()},
			)
		}

		_, err := m.Collection.InsertOne(m.Ctx, inserts)
		if err != nil {
			return err
		}
	} else {
		id := m.docs[idIndex].Value
		updates := append(m.docs[:idIndex], m.docs[idIndex+1:]...)

		if m.option.Timestamp {
			updates = append(updates, bson.E{Key: "updatedAt", Value: time.Now()})
		}
		_, err := m.Collection.UpdateByID(m.Ctx, id, bson.D{{Key: "$set", Value: updates}})
		if err != nil {
			return err
		}
	}

	err = ExecutePostHook(Save, m, m.docs)
	if err != nil {
		return err
	}
	m.docs = nil
	return nil
}

func (m *Model[M]) Pre(nameStr HookName, hookFnc HookFnc[M], async ...bool) {
	names := strings.Split(string(nameStr), "|")
	if len(async) == 0 {
		async = append(async, false)
	}
	for _, name := range names {
		m.preHooks = append(m.preHooks, Hook[M]{
			Name:  HookName(name),
			Func:  hookFnc,
			Async: async[0],
		})
	}
}

func (m *Model[M]) Post(nameStr HookName, hookFnc HookFnc[M], async ...bool) {
	names := strings.Split(string(nameStr), "|")
	if len(async) == 0 {
		async = append(async, false)
	}
	for _, name := range names {
		m.postHooks = append(m.postHooks, Hook[M]{
			Name:  HookName(name),
			Func:  hookFnc,
			Async: async[0],
		})
	}
}

// ToDoc converts an interface{} to a bson.D, suitable for use with the bson and mongo packages.
// If the input is nil, ToDoc returns an empty bson.D.
// ToDoc returns an error if the input cannot be marshaled.
func ToDoc(v interface{}) (doc *bson.D, err error) {
	if v == nil {
		return &bson.D{}, nil
	}
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func (m *Model[M]) getQueryId(id interface{}) (bson.M, error) {
	var query bson.M

	if m.option.ID {
		switch v := id.(type) {
		case string:
			objId, err := primitive.ObjectIDFromHex(id.(string))
			if err != nil {
				return nil, err
			}
			query = bson.M{"_id": objId}
		case primitive.ObjectID:
			query = bson.M{"_id": id}
		default:
			return nil, fmt.Errorf("not support type %v", v)
		}

	} else {
		query = bson.M{"_id": id}
	}

	return query, nil
}
