package mongoose_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_Connect(t *testing.T) {
	loggerOptions := options.
		Logger().
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	connect := mongoose.New(os.Getenv("MONGO_URI"), "test", loggerOptions)
	err := connect.Ping()
	require.Nil(t, err)
}
