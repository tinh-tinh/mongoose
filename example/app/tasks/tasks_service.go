package tasks

import (
	"github.com/tinh-tinh/mongoose"
	"github.com/tinh-tinh/tinhtinh/core"
)

const TASK_SERVICE core.Provide = "TASK_SERVICE"

type tasksService struct {
	model *mongoose.Model[Task]
}

func (s *tasksService) Create() interface{} {
	result, err := s.model.Create(&Task{
		Name:     "huuhuhu",
		Status:   "true",
		TakeTime: 1,
	})
	if err != nil {
		return err
	}
	return result.InsertedID
}

func (s *tasksService) Find() interface{} {
	return nil
}

func (s *tasksService) FindById(id string) interface{} {
	return nil
}

func (s *tasksService) Update(id string, input interface{}) interface{} {
	return nil
}

func (s *tasksService) Delete(id string) interface{} {
	return nil
}

func NewService(module *core.DynamicModule) *core.DynamicProvider {
	svc := module.NewProvider(&tasksService{
		model: mongoose.InjectModel[Task](module, "tasks"),
	}, TASK_SERVICE)

	return svc
}
