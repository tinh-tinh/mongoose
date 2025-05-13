package mongoose

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Config interface {
	string | *options.ClientOptions
}

type Connect struct {
	Client *mongo.Client
	Ctx    context.Context
	DB     string
}

func New[C Config](cfg C) *Connect {
	// You can use type switch if you need runtime behavior
	var connectOptions *options.ClientOptions
	var cs *connstring.ConnString
	switch v := any(cfg).(type) {
	case string:
		// handle string
		connectOptions = options.Client().ApplyURI(v)
		cs, _ = connstring.ParseAndValidate(v)
	case *options.ClientOptions:
		// handle options
		connectOptions = v
	default:
		panic("config is invalid")
	}

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, connectOptions)
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
