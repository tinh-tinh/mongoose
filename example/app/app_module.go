package app

import (
	"github.com/tinh-tinh/mongoose"
	"github.com/tinh-tinh/mongoose/example/app/tasks"
	"github.com/tinh-tinh/tinhtinh/core"
)

func NewModule() *core.DynamicModule {
	appModule := core.NewModule(core.NewModuleOptions{
		Global: true,
		Imports: []core.Module{
			mongoose.ForRoot("mongodb://localhost:27017/hrms"),
			tasks.NewModule,
		},
	})

	appModule.Use()

	return appModule
}
