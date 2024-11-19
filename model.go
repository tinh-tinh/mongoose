package mongoose

import (
	"context"
	"log"
	"reflect"
	"slices"
	"time"

	"github.com/tinh-tinh/tinhtinh/common"
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
	name       string
	docs       []bson.E
	idx        mongo.IndexModel
	Ctx        context.Context
	Collection *mongo.Collection
}

// NewModel returns a new instance of Model[M] with the given connect and name
// name is the name of the collection in the database
// the returned Model[M] is used to interact with the collection in the database
func NewModel[M any](names ...string) *Model[M] {
	var name string
	if len(names) > 0 {
		name = names[0]
	} else {
		var model M
		name = common.GetStructName(model)
	}
	return &Model[M]{
		name: name,
	}
}

// SetConnect sets the context and collection of the model to the given connect.
// It is used internally by the ForFeature function to set the connect of the model.
// The given connect must be a *Connect.
func (m *Model[M]) SetConnect(connect *Connect) {
	m.Ctx = connect.Ctx
	m.Collection = connect.Client.Database(connect.DB).Collection(m.name)
	if !reflect.ValueOf(m.idx).IsZero() {
		_, err := m.Collection.Indexes().CreateOne(m.Ctx, m.idx)
		if err != nil {
			log.Println(err)
		}
	}
}

// GetName returns the name of the collection in the database
func (m *Model[M]) GetName() string {
	return m.name
}

func (m *Model[M]) Index(idx bson.D, unique bool) {
	indexModel := mongo.IndexModel{
		Keys:    idx,
		Options: options.Index().SetUnique(unique),
	}
	m.idx = indexModel
	// _, err := m.Collection.Indexes().CreateOne(m.Ctx, indexModel)
	// if err != nil {
	// 	log.Println(err)
	// }
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
	for i := 0; i < ctInput.NumField(); i++ {
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
	if len(m.docs) == 0 {
		return nil
	}

	idIndex := slices.IndexFunc(m.docs, func(e bson.E) bool {
		return e.Key == "_id"
	})
	if idIndex == -1 {
		inserts := append(m.docs,
			bson.E{Key: "_id", Value: primitive.NewObjectID()},
			bson.E{Key: "createdAt", Value: time.Now()},
			bson.E{Key: "updatedAt", Value: time.Now()},
		)
		_, err := m.Collection.InsertOne(m.Ctx, inserts)
		if err != nil {
			return err
		}
	} else {
		id := m.docs[idIndex].Value
		updates := append(m.docs[:idIndex], m.docs[idIndex+1:]...)
		updates = append(updates, bson.E{Key: "updatedAt", Value: time.Now()})
		_, err := m.Collection.UpdateByID(m.Ctx, id, bson.D{{Key: "$set", Value: updates}})
		if err != nil {
			return err
		}
	}

	m.docs = nil
	return nil
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
