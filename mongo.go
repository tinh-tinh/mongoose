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
	}
}

func (c *Connect) Ping() error {
	return c.Client.Ping(c.Ctx, nil)
}

type Collection struct {
	DB      string
	Name    string
	connect *Connect
}

func NewCollection(c *Connect, name string) *Collection {
	return &Collection{
		DB:      "doban",
		Name:    name,
		connect: c,
	}
}

func (c *Collection) Create(model interface{}) error {
	_, err := c.connect.Client.Database(c.DB).Collection(c.Name).InsertOne(c.connect.Ctx, model)
	if err != nil {
		return err
	}
	return nil
}
