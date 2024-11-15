package mongoose

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_Model(t *testing.T) {
	type Book struct {
		BaseSchema `bson:"inline"`
		Title      string `bson:"title"`
		Author     string `bson:"author"`
	}

	connect := New(os.Getenv("MONGO_URI"), "test")
	model := NewModel[Book]("books")
	model.SetConnect(connect)

	type UpdateBook struct {
		Title string `bson:"title"`
		Level int    `bson:"level"`
	}
	model.Set(&UpdateBook{Title: "abc", Level: 1})
	err := model.Save()
	require.Nil(t, err)
}

func Test_Update(t *testing.T) {
	type Book struct {
		BaseSchema `bson:"inline"`
		Title      string `bson:"title"`
		Author     string `bson:"author"`
	}

	connect := New(os.Getenv("MONGO_URI"), "test")
	model := NewModel[Book]("books")
	model.SetConnect(connect)

	type UpdateBook struct {
		ID    primitive.ObjectID `bson:"_id"`
		Title string             `bson:"title"`
		Level int                `bson:"level"`
	}
	firstOne, err := model.FindOne(nil)
	require.Nil(t, err)
	if firstOne != nil {
		model.Set(&UpdateBook{
			ID:    firstOne.ID,
			Title: "abc",
			Level: 1,
		})
		err := model.Save()
		require.Nil(t, err)
	}
}

func Test_Recusive(t *testing.T) {
	type Address struct {
		Line    string `bson:"line"`
		State   string `bson:"state"`
		City    string `bson:"city"`
		Country string `bson:"country"`
	}

	type Location struct {
		BaseSchema `bson:"inline"`
		Longitude  float64  `bson:"longitude"`
		Latitude   float64  `bson:"latitude"`
		Address    *Address `bson:"address"`
	}

	connect := New(os.Getenv("MONGO_URI"), "test")
	model := NewModel[Location]("locations")
	model.SetConnect(connect)
	_, err := model.Create(&Location{
		Longitude: 1,
		Latitude:  2,
		Address: &Address{
			Line:    "line",
			State:   "state",
			City:    "city",
			Country: "country",
		},
	})
	require.Nil(t, err)

	err = model.Update(nil, &Location{
		Longitude: 2,
		Latitude:  1,
		Address: &Address{
			Line:    "line",
			State:   "state",
			City:    "city",
			Country: "country",
		},
	})
	require.Nil(t, err)
}
