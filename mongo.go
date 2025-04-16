package mongoose

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
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
func New(uri string, opts ...*options.ClientOptions) *Connect {
	// Parse the URI
	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse MongoDB URI: %v\n", err))
	}

	connectOptions := *options.Client().ApplyURI(uri)

	if len(opts) > 0 {
		mergeOpts := append(opts, &connectOptions)
		connectOptions = *options.MergeClientOptions(mergeOpts...)
	}

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, &connectOptions)
	if err != nil {
		panic(err)
	}

	return &Connect{
		Client: client,
		Ctx:    ctx,
		DB:     cs.Database,
	}
}

// Ping pings the MongoDB server to check if the connection is alive.
// It returns an error if the connection is not alive.
func (c *Connect) Ping() error {
	return c.Client.Ping(c.Ctx, nil)
}

func (c *Connect) SetDB(db string) {
	c.DB = db
}
