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

// New creates a new Connect instance by establishing a connection to a MongoDB
// database using the provided URI. It returns a pointer to the Connect struct,
// which contains the MongoDB client, context, and the database name extracted
// from the URI. If the connection fails, the function will panic.
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

// Ping pings the MongoDB server to check if the connection is alive.
// It returns an error if the connection is not alive.
func (c *Connect) Ping() error {
	return c.Client.Ping(c.Ctx, nil)
}
