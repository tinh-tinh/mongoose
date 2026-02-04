package mongoose_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Department struct {
	BaseSchema `bson:"inline"`
	Name                string `bson:"name"`
}

func (d Department) CollectionName() string {
	return "departments"
}

type Employee struct {
	BaseSchema `bson:"inline"`
	Name                string             `bson:"name"`
	Age                 int                `bson:"age"`
	DepartmentID        primitive.ObjectID `bson:"departmentID"`
	Department          *Department        `bson:"department" ref:"departmentID->departments"`
}

func (e Employee) CollectionName() string {
	return "employees"
}

func TestAggregate(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")

	employeeModel := mongoose.NewModel[Employee]()
	employeeModel.SetConnect(connect)

	departmentModel := mongoose.NewModel[Department]()
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
	require.NotNil(t, employees)
	if len(employees) == 0 {
		assert.NotNil(t, employees[0].Department)
	}
}

func TestFindOne(t *testing.T) {
	connect := mongoose.New(os.Getenv("MONGO_URI"))
	connect.SetDB("test")

	employeeModel := mongoose.NewModel[Employee]()
	employeeModel.SetConnect(connect)

	departmentModel := mongoose.NewModel[Department]()
	departmentModel.SetConnect(connect)

	employees, err := employeeModel.FindOne(nil, mongoose.QueryOptions{
		Ref: []string{"departmentID"},
	})
	require.Nil(t, err)
	require.NotNil(t, employees)
}
