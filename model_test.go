package mongoose

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type Task struct {
	BaseSchema `bson:"inline"`
	Name       string `bson:"name"`
	Status     string `bson:"status"`
}

func Test_Create(t *testing.T) {
	connect := New(os.Getenv("MONGO_URI"))
	model := NewModel[Task](connect, "tasks")
	_, err := model.Create(&Task{
		Name:   "haha",
		Status: "true",
	})
	require.Nil(t, err)
}

func Test_CreateMany(t *testing.T) {
	connect := New(os.Getenv("MONGO_URI"))
	model := NewModel[Task](connect, "tasks")
	_, err := model.CreateMany([]*Task{
		{
			Name:   "huuhuhu",
			Status: "true",
		},
		{
			Name:   "lulu",
			Status: "false",
		},
	})
	require.Nil(t, err)
}

type QueryTask struct {
	Name string `bson:"name"`
}

func Test_Find(t *testing.T) {
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

func Test_UpdateOne(t *testing.T) {
	connect := New(os.Getenv("MONGO_URI"))
	model := NewModel[Task](connect, "tasks")
	firstOne, err := model.FindOne(nil)
	require.Nil(t, err)

	if firstOne != nil {
		err := model.Update(&QueryTask{
			Name: firstOne.Name,
		}, &Task{
			Status: "abc",
		})
		require.Nil(t, err)

		reFirst, err := model.FindOne(&QueryTask{
			Name: firstOne.Name,
		})
		require.Nil(t, err)
		require.Equal(t, "abc", reFirst.Status)
	}
}

func Test_UpdateMany(t *testing.T) {
	connect := New(os.Getenv("MONGO_URI"))
	model := NewModel[Task](connect, "tasks")
	err := model.UpdateMany(&QueryTask{
		Name: "haha",
	}, &Task{
		Status: "abc",
	})
	require.Nil(t, err)
}

func Test_DeleteOne(t *testing.T) {
	connect := New(os.Getenv("MONGO_URI"))
	model := NewModel[Task](connect, "tasks")
	err := model.DeleteOne(&QueryTask{
		Name: "huuhuhu",
	})
	require.Nil(t, err)
}

func Test_DeleteMany(t *testing.T) {
	connect := New(os.Getenv("MONGO_URI"))
	model := NewModel[Task](connect, "tasks")
	err := model.DeleteMany(&QueryTask{
		Name: "lulu",
	})
	require.Nil(t, err)
}
