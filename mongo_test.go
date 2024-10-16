package mongoose

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Connect(t *testing.T) {
	connect := New(os.Getenv("MONGO_URI"))
	err := connect.Ping()
	require.Nil(t, err)
}
