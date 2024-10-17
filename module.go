package mongoose

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/utils"
)

const CONNECT_MONGO core.Provide = "CONNECT_MONGO"

func ForRoot(url string) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		mongooseModule := module.New(core.NewModuleOptions{})
		mongooseModule.NewProvider(core.ProviderOptions{
			Name:  CONNECT_MONGO,
			Value: New(url),
		})
		mongooseModule.Export(CONNECT_MONGO)

		return mongooseModule
	}
}

func ForFeature[M any](name string) core.Module {
	var m M
	structName := utils.GetNameStruct(&m)
	return func(module *core.DynamicModule) *core.DynamicModule {
		mongooseModule := module.New(core.NewModuleOptions{})
		mongooseModule.NewProvider(core.ProviderOptions{
			Name: getModelName(structName),
			Factory: func(param ...interface{}) interface{} {
				connect := param[0].(*Connect)
				return NewModel[M](connect, name)
			},
			Inject: []core.Provide{CONNECT_MONGO},
		})
		mongooseModule.Export(getModelName(structName))

		return mongooseModule
	}
}

func getModelName(name string) core.Provide {
	modelName := "Model_" + name

	return core.Provide(modelName)
}

func InjectConnect(module *core.DynamicModule) *Connect {
	return module.Ref(CONNECT_MONGO).(*Connect)
}

func InjectModel[M any](module *core.DynamicModule) *Model[M] {
	name := utils.GetNameStruct(new(M))
	return module.Ref(getModelName(name)).(*Model[M])
}
