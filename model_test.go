package mongoose

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Model(t *testing.T) {
	type Book struct {
		BaseSchema `bson:"inline"`
		Title      string `bson:"title"`
		Author     string `bson:"author"`
	}

	connect := New(os.Getenv("MONGO_URI"))
	model := NewModel[Book](connect, "books")

	type UpdateBook struct {
		Title string `bson:"title"`
		Level int    `bson:"level"`
	}
	model.Set(&UpdateBook{Title: "abc", Level: 1})
	err := model.Save()
	require.Nil(t, err)
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

	connect := New(os.Getenv("MONGO_URI"))
	model := NewModel[Location](connect, "locations")
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
