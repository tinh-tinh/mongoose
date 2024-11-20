package mongoose_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose"
)

func Test_Connect(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"), "test")
	err := connect.Ping()
	require.Nil(t, err)
}
