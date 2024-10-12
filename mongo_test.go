package mongoose

import (
	"fmt"
	"testing"
)

func Test_Connect(t *testing.T) {
	t.Run("Connect", func(t *testing.T) {
		connect := New("mongodb://127.0.0.1:27017/test")
		err := connect.Ping()
		if err != nil {
			t.Error(err)
		}
		fmt.Print("success")
	})
}
