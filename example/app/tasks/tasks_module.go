
package tasks

import "github.com/tinh-tinh/tinhtinh/core"

func NewModule(module *core.DynamicModule) *core.DynamicModule {
	tasksModule := module.New(core.NewModuleOptions{
		Controllers: []core.Controller{NewController},
		Providers:   []core.Provider{NewService},
	})

	return tasksModule
}
	