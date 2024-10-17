package mongoose

import (
	"context"
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
	Ctx        context.Context
	Collection *mongo.Collection
}

func NewModel[M any](connect *Connect, name string) *Model[M] {
	return &Model[M]{
		Ctx:        connect.Ctx,
		Collection: connect.Client.Database(connect.DB).Collection(name),
	}
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
