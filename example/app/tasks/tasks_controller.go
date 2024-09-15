package tasks

import "github.com/tinh-tinh/tinhtinh/core"

func NewController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("tasks")

	ctrl.Post("/", func(ctx core.Ctx) {
		tasksService := module.Ref(TASK_SERVICE).(*tasksService)
		data := tasksService.Create()
		ctx.JSON(core.Map{"data": data})
	})

	ctrl.Get("/", func(ctx core.Ctx) {
		ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Get("/{id}", func(ctx core.Ctx) {
		ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Put("/{id}", func(ctx core.Ctx) {
		ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Delete("/{id}", func(ctx core.Ctx) {
		ctx.JSON(core.Map{"data": "ok"})
	})

	return ctrl
}
