package mongoose_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Models struct {
	mongoose.BaseSchema `bson:"inline"`
	Title               string `bson:"title"`
	Author              string `bson:"author"`
}

func (b Models) CollectionName() string {
	return "models"
}

func Test_Model(t *testing.T) {

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[Models]()
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

func (l Location) CollectionName() string {
	return "recursive"
}

func Test_Recusive(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[Location]()
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

type User struct {
	mongoose.BaseSchema `bson:"inline"`
	Email               string `bson:"email"`
	Name                string `bson:"name"`
}

func (u User) CollectionName() string {
	return "indexes"
}
func TestIndex(t *testing.T) {
	userModel := mongoose.NewModel[User]()
	userModel.Index(bson.D{{Key: "email", Value: 1}}, options.Index().SetUnique(true))

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	userModel.SetConnect(connect)
}

func Test_ToDoc(t *testing.T) {
	_, err := mongoose.ToDoc("nil")
	assert.NotNil(t, err)
}

type Student struct {
	mongoose.BaseTimestamp `bson:"inline"`
	ID                     int    `bson:"_id"`
	FirstName              string `bson:"firstName" validate:"isAlpha"`
	LastName               string `bson:"lastName" validate:"isAlpha"`
	Email                  string `bson:"email" validate:"isEmail"`
}

func (s Student) CollectionName() string {
	return "students"
}

func Test_Validator(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[Student]()
	model.SetConnect(connect)

	err := model.DeleteMany(nil)
	assert.Nil(t, err)

	_, err = model.Create(&Student{
		FirstName: "12",
		LastName:  "222",
		Email:     "111",
	})
	require.NotNil(t, err)

	_, err = model.Create(&Student{
		ID:        1,
		FirstName: "John",
		LastName:  "Dpe",
		Email:     "john@gmail.com",
	})
	require.Nil(t, err)

	_, err = model.FindOneAndUpdate(nil, &Student{Email: "120"})
	assert.NotNil(t, err)

	_, err = model.FindOneAndReplace(nil, &Student{Email: "120"})
	assert.NotNil(t, err)

	_, err = model.CreateMany([]*Student{
		{FirstName: "2"},
	})
	assert.NotNil(t, err)

	_, err = model.CreateMany([]*Student{
		{
			ID:        2,
			FirstName: "Ricardo",
			LastName:  "Kaka",
			Email:     "kaka@gmail.com",
		},
	})
	require.Nil(t, err)

	err = model.Update(map[string]any{}, &Student{FirstName: "$##$$#"})
	assert.NotNil(t, err)

	err = model.UpdateMany(map[string]any{}, &Student{FirstName: "$##$$#"})
	assert.NotNil(t, err)
}
