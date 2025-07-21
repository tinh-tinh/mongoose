package tenancy_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/mongoose/v2"
	"github.com/tinh-tinh/mongoose/v2/tenancy"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Book struct {
	mongoose.BaseSchema `bson:"inline"`
	Title               string `bson:"title"`
	Author              string `bson:"author"`
}

func (b Book) CollectionName() string {
	return "books"
}

func Test_Module(t *testing.T) {
	bookController := func(module core.Module) core.Controller {
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

	bookModule := func(module core.Module) core.Module {
		bookMod := module.New(core.NewModuleOptions{
			Imports: []core.Modules{
				tenancy.ForFeature(mongoose.NewModel[Book]()),
			},
			Controllers: []core.Controllers{bookController},
		})

		return bookMod
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
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
