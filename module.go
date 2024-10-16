package mongoose

import "github.com/tinh-tinh/tinhtinh/core"

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

func InjectConnect(module *core.DynamicModule) *Connect {
	return module.Ref(CONNECT_MONGO).(*Connect)
}

func InjectModel[M any](module *core.DynamicModule, name string) *Model[M] {
	connect := InjectConnect(module)

	return NewModel[M](connect, name)
}
