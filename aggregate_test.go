package mongoose_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAggregate(t *testing.T) {
	type Department struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
	}
	type Employee struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string             `bson:"name"`
		Age                 int                `bson:"age"`
		DepartmentID        primitive.ObjectID `bson:"departmentID"`
		Department          *Department        `bson:"department" ref:"departmentID->departments"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"), "test")

	employeeModel := mongoose.NewModel[Employee]("employees")
	employeeModel.SetConnect(connect)

	departmentModel := mongoose.NewModel[Department]("departments")
	departmentModel.SetConnect(connect)

	_, err := departmentModel.Create(&Department{
		Name: "Finance",
	})
	require.Nil(t, err)

	department, err := departmentModel.FindOne(nil)
	require.Nil(t, err)

	_, err = employeeModel.Create(&Employee{
		Name:         "Kafka",
		Age:          18,
		DepartmentID: department.ID,
	})
	require.Nil(t, err)

	employees, err := employeeModel.Find(nil, mongoose.QueriesOptions{
		Ref: []string{"departmentID"},
	})
	require.Nil(t, err)
	for _, emp := range employees {
		fmt.Println(emp.Department)
	}
}

func TestFindOne(t *testing.T) {
	type Department struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string `bson:"name"`
	}
	type Employee struct {
		mongoose.BaseSchema `bson:"inline"`
		Name                string             `bson:"name"`
		Age                 int                `bson:"age"`
		DepartmentID        primitive.ObjectID `bson:"departmentID"`
		Department          *Department        `bson:"department" ref:"departmentID->departments"`
	}

	connect := mongoose.New(os.Getenv("MONGO_URI"), "test")

	employeeModel := mongoose.NewModel[Employee]("employees")
	employeeModel.SetConnect(connect)

	departmentModel := mongoose.NewModel[Department]("departments")
	departmentModel.SetConnect(connect)

	employees, err := employeeModel.FindOne(nil, mongoose.QueryOptions{
		Ref: []string{"departmentID"},
	})
	require.Nil(t, err)
	require.NotNil(t, employees)
}
