package mongoose

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connect struct {
	Client *mongo.Client
	Ctx    context.Context
	DB     string
}

func New(url string) *Connect {
	opt := options.Client().ApplyURI(url)
	ctx := context.TODO()

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		panic(err)
	}

	return &Connect{
		Client: client,
		Ctx:    ctx,
		DB:     "doban",
	}
}

func (c *Connect) Ping() error {
	return c.Client.Ping(c.Ctx, nil)
}
