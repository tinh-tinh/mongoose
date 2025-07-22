package mongoose_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/mongoose/v2"
)

type MutationTask struct {
	mongoose.BaseSchema `bson:"inline"`
	Name                string `bson:"name"`
	Status              string `bson:"status"`
}

func (t MutationTask) CollectionName() string {
	return "mutations"
}

func Test_Mutation(t *testing.T) {

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[MutationTask]()
	model.SetConnect(connect)

	// Clear before test
	err := model.DeleteMany(nil)
	assert.Nil(t, err)

	// TestCreate
	_, err = model.Create(&MutationTask{
		Name:   "1",
		Status: "true",
	})
	assert.Nil(t, err)

	// TestCreateMany
	_, err = model.CreateMany([]*MutationTask{
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
		}, &MutationTask{
			Status: "abc",
		})
		assert.Nil(t, err)

		reFirst, err := model.FindOne(&QueryTask{
			Name: firstOne.Name,
		})
		assert.Nil(t, err)
		assert.Equal(t, "abc", reFirst.Status)

		err = model.UpdateByID(reFirst.ID, &MutationTask{
			Status: "mno",
		})
		assert.Nil(t, err)
		reFirst, err = model.FindOne(&QueryTask{
			Name: firstOne.Name,
		})
		assert.Nil(t, err)
		assert.Equal(t, "mno", reFirst.Status)

		err = model.UpdateByID("true", &MutationTask{
			Status: "mno",
		})
		assert.NotNil(t, err)
	}

	// TestUpdateMany
	err = model.UpdateMany(&QueryTask{
		Name: "3",
	}, &MutationTask{
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

type Failed struct {
	mongoose.BaseSchema `bson:"inline"`
	Name                string `bson:"name"`
}

func (f Failed) CollectionName() string {
	return "faileds"
}

func Test_Fail(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[Failed]()
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
