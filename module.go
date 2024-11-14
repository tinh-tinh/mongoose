package mongoose

import (
	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/core"
)

const CONNECT_MONGO core.Provide = "CONNECT_MONGO"

// ForRoot creates a module which provides a mongodb connection from a given url.
// The connection is exported as CONNECT_MONGO.
func ForRoot(url string, db string) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		mongooseModule := module.New(core.NewModuleOptions{})
		mongooseModule.NewProvider(core.ProviderOptions{
			Name:  CONNECT_MONGO,
			Value: New(url, db),
		})
		mongooseModule.Export(CONNECT_MONGO)

		return mongooseModule
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
func InjectConnect(module *core.DynamicModule) *Connect {
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
func InjectModel[M any](module *core.DynamicModule, name ...string) *Model[M] {
	var m M
	var modelName string
	if len(name) > 0 {
		modelName = name[0]
	} else {
		modelName = common.GetStructName(&m)
	}
	data, ok := module.Ref(GetModelName(modelName)).(*Model[M])
	if !ok || data == nil {
		model := NewModel[M](InjectConnect(module), modelName)
		module.NewProvider(core.ProviderOptions{
			Name:  GetModelName(modelName),
			Value: model,
		})
		return model
	}

	return data
}
