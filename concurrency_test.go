package mongoose_test

import (
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/mongoose/v2"
)

type ValidationTask struct {
	BaseSchema `bson:"inline"`
	Name       string `bson:"name" validate:"required"`
	Email      string `bson:"email" validate:"isEmail"`
}

func (v ValidationTask) CollectionName() string {
	return "validation_tasks"
}

func Test_Concurrency_Validator(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model := mongoose.NewModel[ValidationTask]()
	model.SetConnect(connect)

	// Clear before test
	err := model.DeleteMany(nil)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	routines := 50
	var successCount int64

	wg.Add(routines)
	for i := 0; i < routines; i++ {
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				// Valid data
				_, err := model.Create(&ValidationTask{
					Name:  "Valid User",
					Email: "abc@gmail.com",
				})
				assert.Nil(t, err)
				atomic.AddInt64(&successCount, 1)
			} else {
				// Invalid data (Email is not valid)
				_, err := model.Create(&ValidationTask{
					Name:  "Kid User",
					Email: "abc",
				})
				assert.NotNil(t, err) // Should fail validation
			}
		}(i)
	}
	wg.Wait()

	// Check count
	count, err := model.Count(nil)
	assert.Nil(t, err)
	assert.Equal(t, successCount, count)
}
