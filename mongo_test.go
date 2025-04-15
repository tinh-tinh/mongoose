package mongoose_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose/v2"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_Connect(t *testing.T) {
	loggerOptions := options.
		Logger().
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	connect := mongoose.New("mongodb://localhost:27017/?replicaSet=rs0", "test", loggerOptions)
	err := connect.Ping()
	require.Nil(t, err)
}
