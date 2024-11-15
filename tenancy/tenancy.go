package tenancy

import (
	"fmt"
	"net/http"

	"github.com/tinh-tinh/mongoose"
	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/core"
)

type Options struct {
	Uri         string
	GetTenantID func(r *http.Request) string
}

const (
	CONNECT_MAPPER  core.Provide = "CONNECT_MAPPER"
	CONNECT_TENANCY core.Provide = "CONNECT_TENANCY"
)

type ConnectMapper map[string]*mongoose.Connect

func CreateConnectMapper(module *core.DynamicModule) *core.DynamicProvider {
	prd := module.NewProvider(core.ProviderOptions{
		Name:  CONNECT_MAPPER,
		Value: make(ConnectMapper),
	})

	return prd
}

func ForRoot(opt Options) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		tenancyModule := module.New(core.NewModuleOptions{})

		CreateConnectMapper(tenancyModule)
		fmt.Println(tenancyModule.DataProviders)

		tenancyModule.NewProvider(core.ProviderOptions{
			Scope: core.Request,
			Name:  CONNECT_TENANCY,
			Factory: func(param ...interface{}) interface{} {
				req := param[0].(*http.Request)
				tenantId := opt.GetTenantID(req)

				connectMapper, ok := param[1].(ConnectMapper)

				if connectMapper == nil || !ok {
					connectMapper = make(ConnectMapper)
				}
				if connectMapper[tenantId] == nil {
					connectMapper[tenantId] = mongoose.New(opt.Uri, tenantId)
				}

				return connectMapper[tenantId]
			},
			Inject: []core.Provide{core.REQUEST, CONNECT_MAPPER},
		})
		tenancyModule.Export(CONNECT_TENANCY)
		return tenancyModule
	}
}

func ForFeature(models ...mongoose.ModelCommon) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		modelModule := module.New(core.NewModuleOptions{Scope: core.Global})
		for _, m := range models {
			modelModule.NewProvider(core.ProviderOptions{
				Scope: core.Request,
				Name:  mongoose.GetModelName(m.GetName()),
				Factory: func(param ...interface{}) interface{} {
					connect := param[0].(*mongoose.Connect)
					m.SetConnect(connect)

					return m
				},
				Inject: []core.Provide{CONNECT_TENANCY},
			})
			modelModule.Export(mongoose.GetModelName(m.GetName()))
		}

		return modelModule
	}
}

func InjectModel[M any](module *core.DynamicModule, name ...string) *mongoose.Model[M] {
	var m M
	var modelName string
	if len(name) > 0 {
		modelName = name[0]
	} else {
		modelName = common.GetStructName(&m)
	}
	data, ok := module.Ref(mongoose.GetModelName(modelName)).(*mongoose.Model[M])
	if !ok {
		return nil
	}

	return data
}
