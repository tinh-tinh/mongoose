package mongoose_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Query_Pre_Hook(t *testing.T) {
	type PreTask struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	type QueryTask struct {
		Name string `bson:"name"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[PreTask]("query_hooks")
	model.SetConnect(connect)

	model.Pre("find|count|findOne|findOneAndDelete|findOneAndReplace|findOneAndUpdate", func(params ...any) error {
		return errors.New("failed to query")
	})

	// TestCount
	_, err := model.Count(nil)
	assert.NotNil(t, err)

	// TestFind
	_, err = model.Find(&QueryTask{Name: "2"})
	assert.NotNil(t, err)

	// TestFindOne
	_, err = model.FindOne(&QueryTask{Name: "2"})
	assert.NotNil(t, err)

	// TestFindOneAndUpdate
	_, err = model.FindOneAndUpdate(&QueryTask{Name: "3"}, &PreTask{Status: "abc"})
	assert.NotNil(t, err)

	// TestFindOneAndReplace
	_, err = model.FindOneAndReplace(&QueryTask{Name: "5"}, &PreTask{Status: "mno"})
	assert.NotNil(t, err)

	// TestFindOneAndDelete
	_, err = model.FindOneAndDelete(&QueryTask{Name: "1"})
	assert.NotNil(t, err)
}

func Test_Query_Post_Hook(t *testing.T) {
	type PostTask struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	type QueryTask struct {
		Name string `bson:"name"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[PostTask]("query_hooks")
	model.SetConnect(connect)

	model.Post("find|count|findOne|findOneAndDelete|findOneAndReplace|findOneAndUpdate", func(params ...any) error {
		return errors.New("failed to query")
	})

	// TestCount
	count, err := model.Count(nil)
	assert.NotNil(t, err)

	if count == 0 {
		_, err = model.CreateMany([]*PostTask{
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
		require.Nil(t, err)
	}

	// TestFind
	_, err = model.Find(&QueryTask{Name: "2"})
	assert.NotNil(t, err)

	// TestFindOne
	_, err = model.FindOne(&QueryTask{Name: "2"})
	assert.NotNil(t, err)

	// TestFindOneAndUpdate
	_, err = model.FindOneAndUpdate(&QueryTask{Name: "3"}, &PostTask{Status: "abc"})
	assert.NotNil(t, err)

	// TestFindOneAndReplace
	_, err = model.FindOneAndReplace(&QueryTask{Name: "5"}, &PostTask{Status: "mno"})
	assert.NotNil(t, err)

	// TestFindOneAndDelete
	_, err = model.FindOneAndDelete(&QueryTask{Name: "1"})
	assert.NotNil(t, err)
}

func Test_Mutation_Pre_Hook(t *testing.T) {
	type Task struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[Task]("mutation_hooks")
	model.SetConnect(connect)

	model.Pre("create|createMany|update|updateMany|delete|deleteMany", func(params ...any) error {
		return errors.New("failed to query")
	})
	err := model.Delete(nil)
	assert.NotNil(t, err)

	err = model.DeleteMany(nil)
	assert.NotNil(t, err)

	_, err = model.Create(&Task{
		Name:   "1",
		Status: "true",
	})
	assert.NotNil(t, err)
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
	assert.NotNil(t, err)

	err = model.Update(nil, &Task{
		Status: "abc",
	})
	assert.NotNil(t, err)

	err = model.UpdateMany(nil, &Task{
		Status: "abc",
	})
	assert.NotNil(t, err)
}

func Test_Mutation_Post_Hook(t *testing.T) {
	type Task struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[Task]("mutation_hooks")
	model.SetConnect(connect)

	model.Post("create|createMany|update|updateMany|delete|deleteMany", func(params ...any) error {
		return errors.New("failed to query")
	})
	_, err := model.Create(&Task{
		Name:   "1",
		Status: "true",
	})
	assert.NotNil(t, err)
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
	assert.NotNil(t, err)

	err = model.Update(nil, &Task{
		Status: "abc",
	})
	assert.NotNil(t, err)

	err = model.UpdateMany(nil, &Task{
		Status: "abc",
	})
	assert.NotNil(t, err)

	err = model.Delete(nil)
	assert.NotNil(t, err)

	err = model.DeleteMany(nil)
	assert.NotNil(t, err)
}

func Test_Save_Hook(t *testing.T) {

	type Book struct {
		mongoose.BaseSchema `bson:"inline"`
		Title               string `bson:"title"`
		Author              string `bson:"author"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[Book]("model_hooks")
	model.SetConnect(connect)

	model.Pre(mongoose.Save, func(params ...any) error {
		docs := params[0].([]bson.E)
		if len(docs) == 0 {
			return errors.New("failed to save")
		}
		return nil
	})

	model.Post(mongoose.Save, func(params ...any) error {
		return errors.New("failed to save")
	})

	err := model.DeleteMany(nil)
	assert.Nil(t, err)

	type CreateBook struct {
		Title string `bson:"title"`
		Level int    `bson:"level"`
	}
	err = model.Save()
	require.NotNil(t, err)

	model.Set(&CreateBook{Title: "abc", Level: 1})
	err = model.Save()
	require.NotNil(t, err)
}

func Test_Validate_Hook(t *testing.T) {
	type PostTask struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
		Status              string `bson:"status"`
	}

	type QueryTask struct {
		Name string `bson:"name"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[PostTask]("validate_hooks")
	model.SetConnect(connect)

	model.Pre(mongoose.Validate, func(params ...any) error {
		docs := params[0].(*PostTask)
		if docs.Status == "abc" {
			return errors.New("failed to save")
		}
		return nil
	})

	model.Post(mongoose.Validate, func(params ...any) error {
		return errors.New("failed to save")
	})

	err := model.Update(&QueryTask{
		Name: "0",
	}, &PostTask{
		Status: "abc",
	})
	assert.NotNil(t, err)

	err = model.Update(&QueryTask{
		Name: "0",
	}, &PostTask{
		Status: "mno",
	})
	assert.NotNil(t, err)
}
