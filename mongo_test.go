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

type Task struct {
	Text      string `bson:"text"`
	Completed bool   `bson:"completed"`
}

func Test_Create(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		connect := New("mongodb://localhost:27017")
		collection := NewCollection(connect, "tasks")
		task := &Task{
			Text:      "test 7",
			Completed: false,
		}
		input := Merge(task)
		err := collection.Create(input)
		if err != nil {
			t.Error(err)
		}
		fmt.Print("success")
	})
}
