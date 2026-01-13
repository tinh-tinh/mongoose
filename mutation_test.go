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

// Test_StrictFilters_Mutation tests that StrictFilters blocks dangerous operators in mutation functions
func Test_StrictFilters_Mutation(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[MutationTask](mongoose.ModelOptions{
		StrictFilters: true,
	})
	model.SetConnect(connect)

	maliciousFilter := map[string]interface{}{
		"name": map[string]interface{}{"$ne": ""},
	}

	// Test Update blocks dangerous operators
	err := model.Update(maliciousFilter, &MutationTask{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "$ne")

	// Test UpdateMany blocks dangerous operators
	err = model.UpdateMany(maliciousFilter, &MutationTask{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "$ne")

	// Test Delete blocks dangerous operators
	err = model.Delete(maliciousFilter)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "$ne")

	// Test DeleteMany blocks dangerous operators
	err = model.DeleteMany(maliciousFilter)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "$ne")

	// Test safe struct filter is allowed
	type SafeFilter struct {
		Name string `bson:"name"`
	}
	err = model.Delete(&SafeFilter{Name: "nonexistent"})
	// This should not error due to sanitization
	if err != nil {
		assert.NotContains(t, err.Error(), "dangerous MongoDB operator")
	}
}
