package mongoose

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseSchema struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type Model[M any] struct {
	Ctx        context.Context
	Collection *mongo.Collection
}

func NewModel[M any](connect *Connect, name string) *Model[M] {
	return &Model[M]{
		Ctx:        connect.Ctx,
		Collection: connect.Client.Database(connect.DB).Collection(name),
	}
}

func (m Model[M]) Create(input *M) (*mongo.InsertOneResult, error) {
	// schema := NewSchema(input)
	ct := reflect.ValueOf(input).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		name := field.Type.Name()
		if name == "BaseSchema" {
			ct.FieldByName(name).Set(reflect.ValueOf(BaseSchema{
				ID:        primitive.NewObjectID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}))
		}
	}
	result, err := m.Collection.InsertOne(m.Ctx, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m Model[M]) Find() *Model[M] {
	return &m
}
