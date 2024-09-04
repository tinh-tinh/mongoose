package mongoose

import (
	"context"
	"reflect"
	"time"

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

func (m Model[M]) Find(filter interface{}, opt ...*options.FindOptions) ([]*M, error) {
	var data []*M

	query, err := ToDoc(filter)
	if err != nil {
		return data, err
	}
	cur, err := m.Collection.Find(m.Ctx, query, opt...)
	if err != nil {
		return data, err
	}

	for cur.Next(m.Ctx) {
		var t M
		err := cur.Decode(&t)
		if err != nil {
			return data, err
		}
		data = append(data, &t)
	}

	if err := cur.Err(); err != nil {
		return data, err
	}

	cur.Close(m.Ctx)

	return data, nil
}

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
