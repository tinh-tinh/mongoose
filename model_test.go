package mongoose

import (
	"fmt"
	"testing"
)

type Task struct {
	BaseSchema `bson:"inline"`
	Name       string `bson:"name"`
	Status     string `bson:"status"`
}

func Test_Model(t *testing.T) {
	t.Run("Test Create", func(t *testing.T) {
		connect := New("mongodb://localhost:27017")
		model := NewModel[Task](connect, "tasks")
		_, err := model.Create(&Task{
			Name:   "huuhuhu",
			Status: "true",
		})

		if err != nil {
			t.Error(err)
		}
	})
}

type QueryTask struct {
	Name string `bson:"name"`
}

func Test_toDoc(t *testing.T) {
	t.Run("Test query", func(t *testing.T) {
		connect := New("mongodb://localhost:27017")
		model := NewModel[Task](connect, "tasks")
		data, err := model.Find(&QueryTask{
			Name: "huuhuhu",
		})
		if err != nil {
			t.Error(err)
		}
		fmt.Println(data)
	})
}
