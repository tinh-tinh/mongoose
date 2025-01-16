package mongoose

import (
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

const CONNECT_MONGO core.Provide = "CONNECT_MONGO"

// ForRoot creates a module which provides a mongodb connection from a given url.
// The connection is exported as CONNECT_MONGO.
func ForRoot(url string, db string) core.Modules {
	return func(module core.Module) core.Module {
		mongooseModule := module.New(core.NewModuleOptions{})
		mongooseModule.NewProvider(core.ProviderOptions{
			Name:  CONNECT_MONGO,
			Value: New(url, db),
		})
		mongooseModule.Export(CONNECT_MONGO)

		return mongooseModule
	}
}

// ForFeature creates a module which provides each model in the given list as a provider.
// The provider of each model is created by calling its SetConnect method with the CONNECT_MONGO
// provider. The name of the provider is the same as the name of the collection, but with "Model_"
// prefixed. The providers are exported by the module.
func ForFeature(models ...ModelCommon) core.Modules {
	return func(module core.Module) core.Module {
		modelModule := module.New(core.NewModuleOptions{})

		for _, m := range models {
			modelModule.NewProvider(core.ProviderOptions{
				Name: GetModelName(m.GetName()),
				Factory: func(param ...interface{}) interface{} {
					connect := param[0].(*Connect)
					m.SetConnect(connect)

					return m
				},
				Inject: []core.Provide{CONNECT_MONGO},
			})
			modelModule.Export(GetModelName(m.GetName()))
		}

		return modelModule
	}
}

// GetModelName returns a unique name for a model provider given a struct name.
// The returned name is in the format "Model_<struct_name>".
func GetModelName(name string) core.Provide {
	modelName := "Model_" + name

	return core.Provide(modelName)
}

// InjectConnect injects the CONNECT_MONGO provider and returns its value as a *Connect.
// The CONNECT_MONGO provider is created by the ForRoot function.
func InjectConnect(module core.Module) *Connect {
	data, ok := module.Ref(CONNECT_MONGO).(*Connect)
	if !ok {
		return nil
	}

	return data
}

// InjectModel injects a model provider and returns its value as a *Model[M].
// The model provider is created by the ForFeature function.
// The name of the provider is the same as the name of the struct,
// but with "Model_" prefixed.
func InjectModel[M any](module core.Module, name ...string) *Model[M] {
	var m M
	var modelName string
	if len(name) > 0 {
		modelName = name[0]
	} else {
		modelName = common.GetStructName(&m)
	}
	data, ok := module.Ref(GetModelName(modelName)).(*Model[M])
	if !ok {
		return nil
	}

	return data
}
