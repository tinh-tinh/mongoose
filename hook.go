package mongoose

import (
	"github.com/tinh-tinh/tinhtinh/v2/common"
)

type HookName string

const (
	Find              HookName = "find"
	Validate          HookName = "validate"
	Save              HookName = "save"
	FindOne           HookName = "findOne"
	FindOneAndDelete  HookName = "findOneAndDelete"
	FindOneAndReplace HookName = "findOneAndReplace"
	FindOneAndUpdate  HookName = "findOneAndUpdate"
	Create            HookName = "create"
	CreateMany        HookName = "createMany"
	Delete            HookName = "delete"
	DeleteMany        HookName = "deleteMany"
	Update            HookName = "update"
	UpdateMany        HookName = "updateMany"
	Count             HookName = "count"
)

type HookFnc[M any] func(params ...any) error

type Hook[M any] struct {
	Name  HookName
	Func  HookFnc[M]
	Async bool
}

func ExecutePreHook[M any](hookName HookName, model *Model[M], params ...any) error {
	hooks := common.Filter(model.preHooks, func(h Hook[M]) bool {
		return h.Name == hookName
	})
	for _, hook := range hooks {
		if hook.Async {
			go hook.Func(params...)
		} else {
			err := hook.Func(params...)
			return err
		}
	}
	return nil
}

func ExecutePostHook[M any](hookName HookName, model *Model[M], params ...any) error {
	hooks := common.Filter(model.postHooks, func(h Hook[M]) bool {
		return h.Name == hookName
	})
	for _, hook := range hooks {
		if hook.Async {
			go hook.Func(params...)
		} else {
			err := hook.Func(params...)
			return err
		}
	}
	return nil
}
