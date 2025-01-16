package mongoose_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/mongoose/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Query(t *testing.T) {
	type Task struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	type QueryTask struct {
		Name string `bson:"name"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"), "test")
	model := mongoose.NewModel[Task]("queries")
	model.SetConnect(connect)

	total, err := model.Count(nil)
	assert.Nil(t, err)

	if total != 0 {
		model.DeleteMany(nil)
	}
	model.CreateMany([]*Task{
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

	// TestFind
	data, err := model.Find(&QueryTask{Name: "2"})
	assert.Nil(t, err)
	assert.Greater(t, len(data), 0)

	// TestFindOne
	first, err := model.FindOne(&QueryTask{Name: "2"})
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.Equal(t, "2", first.Name)

	// TestFindByID
	firstID, err := model.FindByID(first.ID.Hex())
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.Equal(t, first, firstID)

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
	found, err := model.FindOneAndUpdate(&QueryTask{Name: "3"}, &Task{Status: "abc"})
	assert.Nil(t, err)
	assert.Equal(t, "3", found.Name)

	reFirst, err := model.FindOne(&QueryTask{Name: "3"})
	assert.Nil(t, err)
	assert.Equal(t, "abc", reFirst.Status)

	// TestFindByIDAndUpdate
	last, err := model.FindOne(nil, mongoose.QueryOptions{Sort: bson.D{{Key: "name", Value: -1}}})
	assert.Nil(t, err)

	if last != nil {
		_, err := model.FindByIDAndUpdate(last.ID.Hex(), &Task{
			Status: "xyz",
		})
		assert.Nil(t, err)

		reLast, err := model.FindByID(last.ID.Hex())
		assert.Nil(t, err)
		assert.Equal(t, "xyz", reLast.Status)
	}

	// TestFindOneAndReplace
	found, err = model.FindOneAndReplace(&QueryTask{Name: "5"}, &Task{Status: "mno"})
	assert.Nil(t, err)
	assert.NotNil(t, found)

	reFind, err := model.FindByID(found.ID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, "mno", reFind.Status)
	assert.Equal(t, "", reFind.Name)

	// TestFindByIDAndReplace
	found, err = model.FindOne(&QueryTask{Name: "4"})
	assert.Nil(t, err)

	updateFound, err := model.FindByIDAndReplace(found.ID.Hex(), &Task{Status: "ghi"})
	assert.Nil(t, err)

	reCheck, err := model.FindByID(updateFound.ID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, "ghi", reCheck.Status)
	assert.Equal(t, "", reCheck.Name)

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

	reCheck, err = model.FindByID(first.ID.Hex())
	assert.Nil(t, err)
	assert.Nil(t, reCheck)
}
