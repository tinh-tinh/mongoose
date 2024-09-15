package mongoose

import (
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connect struct {
	Client *mongo.Client
	Ctx    context.Context
	DB     string
}

// Url by the format: mongodb://username:password@host:port/database
func New(url string) *Connect {
	uri := strings.Split(url, "/")
	opt := options.Client().ApplyURI(url)
	ctx := context.TODO()

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		panic(err)
	}

	log.Println("Connect to database successful")

	return &Connect{
		Client: client,
		Ctx:    ctx,
		DB:     uri[len(uri)-1],
	}
}

func (c *Connect) Ping() error {
	return c.Client.Ping(c.Ctx, nil)
}
