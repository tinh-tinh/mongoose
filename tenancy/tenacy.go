package tenancy

import (
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
		tenancyModule := module.New(core.NewModuleOptions{})
		tenancyModule.NewReqProvider(string(TENANCY), func(ctx core.Ctx) interface{} {
			tenantId := ctx.Headers(opt.HeaderName)
			url := opt.Uri + core.IfSlashPrefixString(tenantId)

			return mongoose.New(url)
		})
		return tenancyModule
	}
}
