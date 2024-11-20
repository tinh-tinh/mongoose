package mongoose_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose"
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

	model := mongoose.NewModel[Order]("orders")
	model.Index(bson.D{{"code", 1}}, true)

	connect := mongoose.New(os.Getenv("MONGO_URI"), "test")
	model.SetConnect(connect)

	// session, err := connect.Client.StartSession()
	// require.Nil(t, err)

	// defer session.EndSession(connect.Ctx)

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	err := model.Transaction(func(session mongo.SessionContext) error {
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

	// result, err := session.WithTransaction(connect.Ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
	// 	model.SetContext(sessionContext)
	// 	model.Create(&Order{
	// 		Code: "ghi",
	// 	})

	// 	err := model.Update(nil, &Order{
	// 		Code: "abc",
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	result, err := model.FindOne(&QueryOrder{
	// 		Code: "mno",
	// 	})
	// 	return result, err
	// }, txnOptions)
	// fmt.Println(err)
	// fmt.Println(result)

	// err = mongo.WithSession(connect.Ctx, session, func(ctx mongo.SessionContext) error {
	// 	if err = session.StartTransaction(); err != nil {
	// 		return err
	// 	}

	// 	result, _ := model.Create(&Order{
	// 		Code: "kakfa",
	// 	})

	// 	err := model.Update(nil, &Order{
	// 		Code: "mno",
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if err = session.CommitTransaction(connect.Ctx); err != nil {
	// 		return err
	// 	}
	// 	fmt.Println(result)
	// 	return nil
	// })

	// if err != nil {
	// 	if err = session.AbortTransaction(connect.Ctx); err != nil {
	// 		return
	// 	}
	// }
}
