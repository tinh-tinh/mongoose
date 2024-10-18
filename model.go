package mongoose

import (
	"context"
	"reflect"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

type Model[M any] struct {
	docs       []bson.E
	Ctx        context.Context
	Collection *mongo.Collection
}

func NewModel[M any](connect *Connect, name string) *Model[M] {
	return &Model[M]{
		Ctx:        connect.Ctx,
		Collection: connect.Client.Database(connect.DB).Collection(name),
	}
}

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
