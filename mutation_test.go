package mongoose_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose"
)

func Test_Create(t *testing.T) {
	type Task struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"), "test")
	model := mongoose.NewModel[Task]("tasks")
	model.SetConnect(connect)
	_, err := model.Create(&Task{
		Name:   "haha",
		Status: "true",
	})
	require.Nil(t, err)
}

func Test_CreateMany(t *testing.T) {
	type Task struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"), "test")
	model := mongoose.NewModel[Task]("tasks")
	model.SetConnect(connect)
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

func Test_UpdateOne(t *testing.T) {
	type Task struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"), "test")
	model := mongoose.NewModel[Task]("tasks")
	model.SetConnect(connect)
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
	type Task struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"), "test")
	model := mongoose.NewModel[Task]("tasks")
	model.SetConnect(connect)
	err := model.UpdateMany(&QueryTask{
		Name: "haha",
	}, &Task{
		Status: "abc",
	})
	require.Nil(t, err)
}

func Test_DeleteOne(t *testing.T) {
	type Task struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"), "test")
	model := mongoose.NewModel[Task]("tasks")
	model.SetConnect(connect)
	err := model.Delete(&QueryTask{
		Name: "huuhuhu",
	})
	require.Nil(t, err)
}

func Test_DeleteMany(t *testing.T) {
	type Task struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"), "test")
	model := mongoose.NewModel[Task]("tasks")
	model.SetConnect(connect)
	err := model.DeleteMany(&QueryTask{
		Name: "lulu",
	})
	require.Nil(t, err)
}
