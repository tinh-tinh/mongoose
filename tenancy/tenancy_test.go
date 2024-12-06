package tenancy_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose"
	"github.com/tinh-tinh/mongoose/tenancy"
	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Module(t *testing.T) {
	type Book struct {
		mongoose.BaseSchema `bson:"inline"`
		Title               string `bson:"title"`
		Author              string `bson:"author"`
	}

	bookController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("books")

		ctrl.Post("", func(ctx core.Ctx) error {
			service := tenancy.InjectModel[Book](module, ctx)
			if service == nil {
				return common.InternalServerException(ctx.Res(), "service is nil")
			}
			data, err := service.Create(&Book{
				Title:  "The Catcher in the Rye",
				Author: "J. D. Salinger",
			})

			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}

			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		return ctrl
	}

	bookModule := func(module *core.DynamicModule) *core.DynamicModule {
		bookMod := module.New(core.NewModuleOptions{
			Imports: []core.Module{
				tenancy.ForFeature(mongoose.NewModel[Book]()),
			},
			Controllers: []core.Controller{bookController},
		})

		return bookMod
	}

	appModule := func() *core.DynamicModule {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Module{
				tenancy.ForRoot(tenancy.Options{
					GetTenantID: func(r *http.Request) string {
						return r.Header.Get("x-tenant-id")
					},
					Uri: os.Getenv("MONGO_URI"),
				}),
				bookModule,
			},
		})

		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/app")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	req, err := http.NewRequest("POST", testServer.URL+"/app/books", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "anc")

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest("POST", testServer.URL+"/app/books", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "xyz")

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
