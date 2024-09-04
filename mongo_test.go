package mongoose

import (
	"fmt"
	"testing"
)

func Test_Connect(t *testing.T) {
	t.Run("Connect", func(t *testing.T) {
		connect := New("mongodb://localhost:27017")
		err := connect.Ping()
		if err != nil {
			t.Error(err)
		}
		fmt.Print("success")
	})
}
