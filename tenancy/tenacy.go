package tenancy

import (
	"net/http"

	"github.com/tinh-tinh/mongoose"
	"github.com/tinh-tinh/tinhtinh/core"
)

const TENANCY core.Provide = "TENANCY"

type Options struct {
	Uri        string
	HeaderName string
}

func ForRoot(opt Options) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		tenancyModule := module.New(core.NewModuleOptions{
			Scope: core.Request,
		})
		tenancyModule.NewProvider(core.ProviderOptions{
			Name: TENANCY,
			Factory: func(param ...interface{}) interface{} {
				req := param[0].(*http.Request)
				tenantId := req.Header.Get(opt.HeaderName)

				url := opt.Uri + core.IfSlashPrefixString(tenantId)
				return mongoose.New(url)
			},
			Inject: []core.Provide{core.REQUEST},
		})
		tenancyModule.Export(TENANCY)
		return tenancyModule
	}
}
