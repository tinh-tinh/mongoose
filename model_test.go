package mongoose

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
			BaseSchema: BaseSchema{
				ID:        primitive.NewObjectID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:   "abv",
			Status: "true",
		})

		if err != nil {
			t.Error(err)
		}
	})
}
