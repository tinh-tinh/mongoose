package mongoose

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

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
	result, err := m.Collection.InsertOne(m.Ctx, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m Model[M]) Find() *Model[M] {
	return &m
}
