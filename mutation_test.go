package mongoose_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/mongoose/v2"
)

func Test_Mutation(t *testing.T) {
	type Task struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[Task]("mutations")
	model.SetConnect(connect)

	// Clear before test
	err := model.DeleteMany(nil)
	assert.Nil(t, err)

	// TestCreate
	_, err = model.Create(&Task{
		Name:   "1",
		Status: "true",
	})
	assert.Nil(t, err)

	// TestCreateMany
	_, err = model.CreateMany([]*Task{
		{
			Name:   "2",
			Status: "true",
		},
		{
			Name:   "3",
			Status: "false",
		},
		{
			Name:   "3",
			Status: "true",
		},
	})
	assert.Nil(t, err)

	type QueryTask struct {
		Name string `bson:"name"`
	}
	firstOne, err := model.FindOne(nil)
	assert.Nil(t, err)

	// Test Update
	if firstOne != nil {
		err := model.Update(&QueryTask{
			Name: firstOne.Name,
		}, &Task{
			Status: "abc",
		})
		assert.Nil(t, err)

		reFirst, err := model.FindOne(&QueryTask{
			Name: firstOne.Name,
		})
		assert.Nil(t, err)
		assert.Equal(t, "abc", reFirst.Status)

		err = model.UpdateByID(reFirst.ID, &Task{
			Status: "mno",
		})
		assert.Nil(t, err)
		reFirst, err = model.FindOne(&QueryTask{
			Name: firstOne.Name,
		})
		assert.Nil(t, err)
		assert.Equal(t, "mno", reFirst.Status)

		err = model.UpdateByID("true", &Task{
			Status: "mno",
		})
		assert.NotNil(t, err)
	}

	// TestUpdateMany
	err = model.UpdateMany(&QueryTask{
		Name: "3",
	}, &Task{
		Status: "abc",
	})
	assert.Nil(t, err)

	reCheck, err := model.FindOne(&QueryTask{Name: "3"})
	assert.Nil(t, err)
	assert.Equal(t, "abc", reCheck.Status)

	// TestDelete
	err = model.Delete(&QueryTask{Name: "2"})
	assert.Nil(t, err)

	reCheck, err = model.FindOne(&QueryTask{Name: "2"})
	assert.Nil(t, err)
	assert.Nil(t, reCheck)

	firstOne, err = model.FindOne(nil)
	assert.Nil(t, err)
	if firstOne != nil {
		err = model.DeleteByID(firstOne.ID)
		assert.Nil(t, err)
	}

	err = model.DeleteByID(true)
	assert.NotNil(t, err)

	// TestDeleteMany
	err = model.DeleteMany(&QueryTask{Name: "3"})
	assert.Nil(t, err)

	reCheck, err = model.FindOne(&QueryTask{Name: "3"})
	assert.Nil(t, err)
	assert.Nil(t, reCheck)
}

func Test_Fail(t *testing.T) {
	type Failed struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[Failed]("faileds")
	model.SetConnect(connect)

	err := model.Update("abc", nil)
	assert.NotNil(t, err)

	err = model.UpdateMany("abc", nil)
	assert.NotNil(t, err)

	err = model.Delete("abc")
	assert.NotNil(t, err)

	err = model.DeleteMany("abc")
	assert.NotNil(t, err)
}
