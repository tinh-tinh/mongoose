package tenancy

import (
	"net/http"
	"reflect"
	"sync"

	"github.com/tinh-tinh/mongoose/v2"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
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

// connectMapperMu protects concurrent access to ConnectMapper
var connectMapperMu sync.RWMutex

// CreateConnectMapper creates a provider named CONNECT_MAPPER which is a map of
// tenant_id to *mongoose.Connect. The map is used to store the connection to the
// database for each tenant. The provider is created by the ForRoot function and
// is used by the ForFeature function to inject the connection of the tenant to the
// model.
func CreateConnectMapper(module core.Module) core.Provider {
	prd := module.NewProvider(core.ProviderOptions{
		Name:  CONNECT_MAPPER,
		Value: make(ConnectMapper),
	})

	return prd
}

// ForRoot creates a module that manages tenant-specific MongoDB connections.
// It utilizes the provided Options to extract the tenant ID from each HTTP request
// and maintain connections in a ConnectMapper. The function creates the CONNECT_MAPPER
// provider to store these connections, and the CONNECT_TENANCY provider to inject
// tenant-specific connections into the models. This setup allows each tenant to have
// a dedicated MongoDB connection based on their tenant ID.

func ForRoot(opt Options) core.Modules {
	return func(module core.Module) core.Module {
		tenancyModule := module.New(core.NewModuleOptions{})

		CreateConnectMapper(tenancyModule)
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

				// Check with read lock first
				connectMapperMu.RLock()
				conn := connectMapper[tenantId]
				connectMapperMu.RUnlock()

				if conn == nil {
					// Acquire write lock for initialization
					connectMapperMu.Lock()
					defer connectMapperMu.Unlock()
					// Double-check after acquiring write lock
					if connectMapper[tenantId] == nil {
						connectMapper[tenantId] = mongoose.New(opt.Uri)
						connectMapper[tenantId].SetDB(tenantId)
					}
					conn = connectMapper[tenantId]
				}

				return conn
			},
			Inject: []core.Provide{core.REQUEST, CONNECT_MAPPER},
		})
		tenancyModule.Export(CONNECT_TENANCY)
		return tenancyModule
	}
}

// ForFeature creates a module which provides each model in the given list as a provider.
// The provider of each model is created by calling its SetConnect method with the CONNECT_TENANCY
// provider. The name of the provider is the same as the name of the collection, but with "Model_"
// prefixed. The providers are exported by the module.
func ForFeature(models ...mongoose.ModelCommon) core.Modules {
	return func(module core.Module) core.Module {
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

// InjectModel injects a model provider and returns its value as a *Model[M].
// The model provider is created by the ForFeature function.
// The name of the provider is the same as the name of the struct,
// but with "Model_" prefixed.
func InjectModel[M any](module core.Module, ctx core.Ctx, name ...string) *mongoose.Model[M] {
	var model M
	ctModel := reflect.ValueOf(&model).Elem()

	fnc := ctModel.MethodByName("CollectionName")

	var modelName string
	if fnc.IsValid() {
		modelName = fnc.Call(nil)[0].String()
	} else {
		modelName = common.GetStructName(model)
	}
	data, ok := module.Ref(mongoose.GetModelName(modelName), ctx).(*mongoose.Model[M])
	if !ok {
		return nil
	}

	return data
}
