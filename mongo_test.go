package mongoose_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose/v2"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_Connect(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	err := connect.Ping()
	require.Nil(t, err)
}

func Test_ConnectFail(t *testing.T) {
	require.Panics(t, func() {
		connect := mongoose.New("http://localhost:27017")
		require.Nil(t, connect)
	})
}

func Test_ConnecPanic(t *testing.T) {
	require.Panics(t, func() {
		connect := mongoose.New(mongoose.Options{})
		require.Nil(t, connect)
	})
}

func Test_Retry(t *testing.T) {
	require.Panics(t, func() {
		connect := mongoose.New(mongoose.Options{
			ClientOptions: options.Client().ApplyURI("MONGO_URI"),
			RetryOptions: mongoose.RetryOptions{
				Retry: 3,
				Delay: 1 * time.Second, // 1 second
			},
		})
		require.Nil(t, connect)
	})
}
