package mongoose_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_Model(t *testing.T) {
	type Book struct {
		mongoose.BaseSchema `bson:"inline"`
		Title               string `bson:"title"`
		Author              string `bson:"author"`
	}

	connect := mongoose.New("mongodb://localhost:27017/?replicaSet=rs0", "test")
	model := mongoose.NewModel[Book]("models")
	model.SetConnect(connect)

	err := model.DeleteMany(nil)
	assert.Nil(t, err)

	type CreateBook struct {
		Title string `bson:"title"`
		Level int    `bson:"level"`
	}
	model.Set(&CreateBook{Title: "abc", Level: 1})
	err = model.Save()
	require.Nil(t, err)

	firstOne, err := model.FindOne(nil)
	require.Nil(t, err)
	require.NotNil(t, firstOne)
	require.Equal(t, "abc", firstOne.Title)

	type UpdateBook struct {
		ID    primitive.ObjectID `bson:"_id"`
		Title string             `bson:"title"`
		Level int                `bson:"level"`
	}

	model.Set(&UpdateBook{ID: firstOne.ID, Title: "xyz", Level: 1})
	err = model.Save()
	require.Nil(t, err)

	reFirst, err := model.FindOne(nil)
	require.Nil(t, err)
	require.Equal(t, "xyz", reFirst.Title)
}

func Test_Recusive(t *testing.T) {
	type Address struct {
		Line    string `bson:"line"`
		State   string `bson:"state"`
		City    string `bson:"city"`
		Country string `bson:"country"`
	}

	type Location struct {
		mongoose.BaseSchema `bson:"inline"`
		Longitude           float64  `bson:"longitude"`
		Latitude            float64  `bson:"latitude"`
		Address             *Address `bson:"address"`
	}

	connect := mongoose.New("mongodb://localhost:27017/?replicaSet=rs0", "test")
	model := mongoose.NewModel[Location]("recursive")
	model.SetConnect(connect)

	err := model.DeleteMany(nil)
	assert.Nil(t, err)

	_, err = model.Create(&Location{
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
	firstOne, err := model.FindOne(nil)

	require.Nil(t, err)
	require.NotNil(t, firstOne)

	require.Equal(t, float64(1), firstOne.Longitude)
	require.Equal(t, float64(2), firstOne.Latitude)
	require.Equal(t, "line", firstOne.Address.Line)
	require.Equal(t, "state", firstOne.Address.State)
	require.Equal(t, "city", firstOne.Address.City)
	require.Equal(t, "country", firstOne.Address.Country)

	err = model.Update(nil, &Location{
		Longitude: 2,
		Latitude:  1,
		Address: &Address{
			Line:    "line2",
			State:   "state2",
			City:    "city2",
			Country: "country2",
		},
	})
	require.Nil(t, err)

	reFirst, err := model.FindOne(nil)
	require.Nil(t, err)
	require.Equal(t, float64(2), reFirst.Longitude)
	require.Equal(t, float64(1), reFirst.Latitude)
	require.Equal(t, "line2", reFirst.Address.Line)
	require.Equal(t, "state2", reFirst.Address.State)
	require.Equal(t, "city2", reFirst.Address.City)
	require.Equal(t, "country2", reFirst.Address.Country)
}

func TestIndex(t *testing.T) {
	type User struct {
		mongoose.BaseSchema `bson:"inline"`
		Email               string `bson:"email"`
		Name                string `bson:"name"`
	}
	userModel := mongoose.NewModel[User]("indexes")
	userModel.Index(bson.D{{Key: "email", Value: 1}}, true)

	connect := mongoose.New("mongodb://localhost:27017/?replicaSet=rs0", "test")
	userModel.SetConnect(connect)
}
