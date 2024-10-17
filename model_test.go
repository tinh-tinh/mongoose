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
