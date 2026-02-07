package mongoose_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/mongoose/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type QueryTask struct {
	BaseSchema `bson:"inline"`
	Name                string `bson:"name"`
	Status              string `bson:"status"`
}

func (t QueryTask) CollectionName() string {
	return "queries"
}

func Test_Query(t *testing.T) {
	type FindParamTask struct {
		Name string `bson:"name"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[QueryTask]()
	model.SetConnect(connect)

	total, err := model.Count(nil)
	assert.Nil(t, err)

	if total != 0 {
		err := model.DeleteMany(nil)
		assert.Nil(t, err)
	}
	_, err = model.CreateMany([]*QueryTask{
		{
			Name:   "1",
			Status: "true",
		},
		{
			Name:   "2",
			Status: "true",
		},
		{
			Name:   "3",
			Status: "false",
		},
		{
			Name:   "4",
			Status: "true",
		},
		{
			Name:   "2",
			Status: "true",
		},
		{
			Name:   "5",
			Status: "false",
		},
		{
			Name:   "6",
			Status: "true",
		},
	})
	assert.Nil(t, err)

	// TestFind
	data, err := model.Find(&FindParamTask{Name: "2"})
	assert.Nil(t, err)
	assert.Greater(t, len(data), 0)

	// TestFindOne
	first, err := model.FindOne(&FindParamTask{Name: "2"})
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.Equal(t, "2", first.Name)

	// TestFindByID
	firstID, err := model.FindByID(first.ID.Hex())
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.Equal(t, first, firstID)

	_, err = model.FindByID(true)
	assert.NotNil(t, err)

	// TestFindOptions
	data, err = model.Find(nil, mongoose.QueriesOptions{
		Sort:       bson.D{{Key: "name", Value: 1}},
		Projection: bson.D{{Key: "name", Value: 1}},
		Skip:       1,
		Limit:      2,
	})
	assert.Nil(t, err)
	assert.Len(t, data, 2)
	assert.Equal(t, "2", data[0].Name)
	assert.Equal(t, "2", data[1].Name)
	assert.Equal(t, "", data[0].Status)

	// TestFindOneOptions
	firstOne, err := model.FindOne(nil, mongoose.QueryOptions{
		Sort:       bson.D{{Key: "name", Value: 1}},
		Projection: bson.D{{Key: "name", Value: 1}},
	})
	assert.Nil(t, err)
	assert.Equal(t, "1", firstOne.Name)
	assert.Equal(t, "", firstOne.Status)

	// TestFindOneAndUpdate
	found, err := model.FindOneAndUpdate(&FindParamTask{Name: "3"}, &QueryTask{Status: "abc"})
	assert.Nil(t, err)
	assert.Equal(t, "3", found.Name)

	reFirst, err := model.FindOne(&FindParamTask{Name: "3"})
	assert.Nil(t, err)
	assert.Equal(t, "abc", reFirst.Status)

	// TestFindByIDAndUpdate
	last, err := model.FindOne(nil, mongoose.QueryOptions{Sort: bson.D{{Key: "name", Value: -1}}})
	assert.Nil(t, err)

	if last != nil {
		_, err := model.FindByIDAndUpdate(last.ID.Hex(), &QueryTask{
			Status: "xyz",
		})
		assert.Nil(t, err)

		reLast, err := model.FindByID(last.ID.Hex())
		assert.Nil(t, err)
		assert.Equal(t, "xyz", reLast.Status)

		_, err = model.FindByIDAndUpdate(true, &QueryTask{
			Status: "xyz",
		})
		assert.NotNil(t, err)
	}

	// TestFindOneAndReplace
	found, err = model.FindOneAndReplace(&FindParamTask{Name: "5"}, &QueryTask{Status: "mno"})
	assert.Nil(t, err)
	assert.NotNil(t, found)

	reFind, err := model.FindByID(found.ID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, "mno", reFind.Status)
	assert.Equal(t, "", reFind.Name)

	// TestFindByIDAndReplace
	found, err = model.FindOne(&QueryTask{Name: "4"})
	assert.Nil(t, err)
	if found != nil {
		updateFound, err := model.FindByIDAndReplace(found.ID.Hex(), &QueryTask{Status: "ghi"})
		assert.Nil(t, err)

		reCheck, err := model.FindByID(updateFound.ID.Hex())
		assert.Nil(t, err)
		assert.Equal(t, "ghi", reCheck.Status)
		assert.Equal(t, "", reCheck.Name)
	}

	_, err = model.FindByIDAndReplace(true, &QueryTask{
		Status: "xyz",
	})
	assert.NotNil(t, err)
	// TestFindOneAndDelete
	_, err = model.FindOneAndDelete(&QueryTask{Name: "1"})
	assert.Nil(t, err)

	found, err = model.FindOne(&QueryTask{Name: "1"})
	assert.Nil(t, err)
	assert.Nil(t, found)

	first, err = model.FindOne(nil)
	assert.Nil(t, err)

	_, err = model.FindByIDAndDelete(first.ID.Hex())
	assert.Nil(t, err)

	reCheck, err := model.FindByID(first.ID.Hex())
	assert.Nil(t, err)
	assert.Nil(t, reCheck)

	_, err = model.FindByIDAndDelete(true)
	assert.NotNil(t, err)
}

type SpecialTask struct {
	ID     int    `bson:"_id"`
	Name   string `bson:"name"`
	Status string `bson:"status"`
}

func (s SpecialTask) CollectionName() string {
	return "sp_tasks"
}

func Test_NotTimestamp(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[SpecialTask](mongoose.ModelOptions{
		Timestamp: false,
		ID:        false,
	})
	model.SetConnect(connect)

	total, err := model.Count(nil)
	assert.Nil(t, err)

	if total == 0 {
		_, err = model.Create(&SpecialTask{
			ID:     1,
			Name:   "abv",
			Status: "true",
		})
		assert.Nil(t, err)
	}

	data, err := model.FindByID(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, data.ID)

	data, err = model.FindByIDAndUpdate(1, &SpecialTask{
		Status: "xtz",
	})
	assert.Nil(t, err)

	reLast, err := model.FindByID(data.ID)
	assert.Nil(t, err)
	assert.Equal(t, "xtz", reLast.Status)

	updateFound, err := model.FindByIDAndReplace(data.ID, &SpecialTask{Status: "ghi"})
	assert.Nil(t, err)

	reCheck, err := model.FindByID(updateFound.ID)
	assert.Nil(t, err)
	assert.Equal(t, "ghi", reCheck.Status)
	assert.Equal(t, "", reCheck.Name)

	_, err = model.FindByIDAndDelete(reCheck.ID)
	assert.Nil(t, err)

	reCheck, err = model.FindByID(reCheck.ID)
	assert.Nil(t, err)
	assert.Nil(t, reCheck)
}

func Test_FailedFind(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[Failed]()
	model.SetConnect(connect)

	_, err := model.FindOne("abc")
	assert.NotNil(t, err)

	_, err = model.Find("abc")
	assert.NotNil(t, err)

	_, err = model.Count("abc")
	assert.NotNil(t, err)

	_, err = model.FindOneAndUpdate("abc", &Failed{})
	assert.NotNil(t, err)

	_, err = model.FindOneAndReplace("abc", &Failed{})
	assert.NotNil(t, err)

	_, err = model.FindOneAndDelete("abc")
	assert.NotNil(t, err)
}

// Test_StrictFilters_Query tests that StrictFilters blocks dangerous operators in query functions
func Test_StrictFilters_Query(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[QueryTask](mongoose.ModelOptions{
		StrictFilters: true,
	})
	model.SetConnect(connect)

	maliciousFilter := map[string]interface{}{
		"name": map[string]interface{}{"$ne": ""},
	}

	// Test FindOne blocks dangerous operators
	_, err := model.FindOne(maliciousFilter)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "$ne")

	// Test Find blocks dangerous operators
	_, err = model.Find(maliciousFilter)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "$ne")

	// Test Count blocks dangerous operators
	_, err = model.Count(maliciousFilter)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "$ne")

	// Test FindOneAndUpdate blocks dangerous operators
	_, err = model.FindOneAndUpdate(maliciousFilter, &QueryTask{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "$ne")

	// Test FindOneAndReplace blocks dangerous operators
	_, err = model.FindOneAndReplace(maliciousFilter, &QueryTask{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "$ne")

	// Test FindOneAndDelete blocks dangerous operators
	_, err = model.FindOneAndDelete(maliciousFilter)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "$ne")

	// Test safe struct filter is allowed
	type SafeFilter struct {
		Name string `bson:"name"`
	}
	_, err = model.FindOne(&SafeFilter{Name: "test"})
	// This should not error due to sanitization (may error due to no docs found, which is fine)
	if err != nil {
		assert.NotContains(t, err.Error(), "dangerous MongoDB operator")
	}
}
