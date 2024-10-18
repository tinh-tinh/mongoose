package mongoose

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/utils"
)

const CONNECT_MONGO core.Provide = "CONNECT_MONGO"

// ForRoot creates a module which provides a mongodb connection from a given url.
// The connection is exported as CONNECT_MONGO.
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

// ForFeature creates a module which provides a model for a given struct.
// The model is exported as a provider with the name of the struct.
// The model is created with the given collection name.
// The module injects the CONNECT_MONGO provider.
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

// getModelName returns a unique name for a model provider given a struct name.
// The returned name is in the format "Model_<struct_name>".
func getModelName(name string) core.Provide {
	modelName := "Model_" + name

	return core.Provide(modelName)
}

// InjectConnect injects the CONNECT_MONGO provider and returns its value as a *Connect.
// The CONNECT_MONGO provider is created by the ForRoot function.
func InjectConnect(module *core.DynamicModule) *Connect {
	return module.Ref(CONNECT_MONGO).(*Connect)
}

// InjectModel injects a model provider and returns its value as a *Model[M].
// The model provider is created by the ForFeature function.
// The name of the provider is the same as the name of the struct,
// but with "Model_" prefixed.
func InjectModel[M any](module *core.DynamicModule) *Model[M] {
	name := utils.GetNameStruct(new(M))
	return module.Ref(getModelName(name)).(*Model[M])
}
