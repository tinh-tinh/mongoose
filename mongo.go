package mongoose

import (
	"context"
	"fmt"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/common/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Config interface {
	string | Options
}

type RetryOptions struct {
	Retry int
	Delay time.Duration // in milliseconds
}

type Options struct {
	*options.ClientOptions
	RetryOptions RetryOptions
}

type Connect struct {
	Client *mongo.Client
	Ctx    context.Context
	DB     string
}

func New[C Config](cfg C) *Connect {
	// You can use type switch if you need runtime behavior
	var connectOptions *options.ClientOptions
	var retryOptions RetryOptions
	var cs *connstring.ConnString
	switch v := any(cfg).(type) {
	case string:
		// handle string
		connectOptions = options.Client().ApplyURI(v)
		cs, _ = connstring.ParseAndValidate(v)
	case Options:
		// handle options
		connectOptions = v.ClientOptions
		retryOptions = v.RetryOptions
	default:
		panic("config is invalid")
	}

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, connectOptions)
	if err != nil {
		if retryOptions.Retry > 0 {
			fmt.Printf("%s %s %s %s\n",
				color.Green("MONGOOSE"),
				color.White("Failed to connect to MongoDB:"),
				color.Red(err.Error()),
				color.Yellow(fmt.Sprintf("Retrying attempt remain %d", retryOptions.Retry)),
			)
			time.Sleep(retryOptions.Delay)
			return New(Options{
				ClientOptions: connectOptions,
				RetryOptions: RetryOptions{
					Retry: retryOptions.Retry - 1,
					Delay: retryOptions.Delay,
				},
			})
		}
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
