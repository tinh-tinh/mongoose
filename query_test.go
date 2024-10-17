package mongoose

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Find(t *testing.T) {
	type Task struct {
		BaseSchema `bson:"inline"`
		Name       string `bson:"name"`
		Status     string `bson:"status"`
	}

	connect := New(os.Getenv("MONGO_URI"))
	model := NewModel[Task](connect, "tasks")
	data, err := model.Find(&QueryTask{
		Name: "haha",
	})
	require.Nil(t, err)
	if len(data) > 0 {
		require.Equal(t, "haha", data[0].Name)
	}
}

func Test_FindOne(t *testing.T) {
	type Task struct {
		BaseSchema `bson:"inline"`
		Name       string `bson:"name"`
		Status     string `bson:"status"`
	}

	connect := New(os.Getenv("MONGO_URI"))
	model := NewModel[Task](connect, "tasks")
	data, err := model.FindOne(&QueryTask{
		Name: "lulu",
	})
	require.Nil(t, err)
	if data != nil {
		require.Equal(t, "lulu", data.Name)
	}

	data, err = model.FindOne(&QueryTask{
		Name: "haha",
	})
	require.Nil(t, err)
	if data != nil {
		require.Equal(t, "haha", data.Name)
	}
}

func Test_FindByID(t *testing.T) {
	type Task struct {
		BaseSchema `bson:"inline"`
		Name       string `bson:"name"`
		Status     string `bson:"status"`
	}

	connect := New(os.Getenv("MONGO_URI"))
	model := NewModel[Task](connect, "tasks")
	firstOne, err := model.FindOne(nil)
	require.Nil(t, err)

	if firstOne != nil {
		data, err := model.FindByID(firstOne.ID.Hex())
		require.Nil(t, err)
		if data != nil {
			require.Equal(t, firstOne.Name, data.Name)
		}
	}
}

func Test_FindOptions(t *testing.T) {
	connect := New(os.Getenv("MONGO_URI"))

	type Course struct {
		BaseSchema `bson:"inline"`
		Title      string `bson:"title"`
		Enrollment int    `bson:"enrollment"`
		CourseId   string `bson:"course_id"`
	}
	model := NewModel[Course](connect, "courses")

	count, err := model.Count(nil)
	require.Nil(t, err)
	if count == 0 {
		_, err := model.CreateMany([]*Course{
			{Title: "World Fiction", CourseId: "PSY2030", Enrollment: 36},
			{Title: "Abstract Algebra", CourseId: "PSY2031", Enrollment: 60},
			{Title: "Modern Poetry", CourseId: "PSY2032", Enrollment: 12},
			{Title: "Plate Tectonics", CourseId: "PSY2033", Enrollment: 35},
		})
		require.Nil(t, err)
	}

	data, err := model.FindOne(nil, QueryOptions{
		Sort: bson.D{{"enrollment", -1}},
	})

	require.Nil(t, err)

	if data != nil {
		require.Equal(t, "Abstract Algebra", data.Title)
		require.Equal(t, int(60), data.Enrollment)
	}

	list, err := model.Find(nil, QueriesOptions{
		Sort:  bson.D{{"enrollment", -1}},
		Limit: 2,
		Skip:  2,
		Projection: bson.D{
			{"title", 1},
			{"enrollment", 1},
		},
	})
	require.Nil(t, err)
	require.Equal(t, 2, len(list))

	require.Equal(t, "Plate Tectonics", list[0].Title)
	require.Equal(t, "Modern Poetry", list[1].Title)

	require.Empty(t, list[0].CourseId)
}
