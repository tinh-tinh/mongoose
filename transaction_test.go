package mongoose_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func TestTransaction(t *testing.T) {
	type Order struct {
		mongoose.BaseSchema `bson:"inline"`
		Code                string `bson:"code"`
		Paid                int    `bson:"paid"`
	}
	type QueryOrder struct {
		Code string `bson:"code"`
	}

	model := mongoose.NewModel[Order]("transactions")
	model.Index(bson.D{{Key: "code", Value: 1}}, true)

	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")
	model.SetConnect(connect)

	err := model.DeleteMany(nil)
	assert.Nil(t, err)

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	err = model.Transaction(func(session mongo.SessionContext) error {
		model.Create(&Order{
			Code: "kai",
		})

		err := model.Update(nil, &Order{
			Code: "vin",
		})
		if err != nil {
			return err
		}
		result, err := model.FindOne(&QueryOrder{
			Code: "vin",
		})
		fmt.Println(result, err)
		return nil
	}, txnOptions)

	require.Nil(t, err)

	err = model.Transaction(func(session mongo.SessionContext) error {
		model.Create(&Order{
			Code: "aul",
		})

		err := model.Update(nil, &Order{
			Code: "aul",
		})
		if err != nil {
			return err
		}
		result, err := model.FindOne(&QueryOrder{
			Code: "vin",
		})
		fmt.Println(result, err)
		return err
	}, txnOptions)

	require.NotNil(t, err)
}
