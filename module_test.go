package mongoose

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Module(t *testing.T) {
	type Book struct {
		BaseSchema `bson:"inline"`
		Title      string `bson:"title"`
		Author     string `bson:"author"`
	}

	bookController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("books")

		ctrl.Get("connect", func(ctx core.Ctx) error {
			connect := InjectConnect(module)
			return ctx.JSON(core.Map{
				"data": connect,
			})
		})

		ctrl.Post("", func(ctx core.Ctx) error {
			service := InjectModel[Book](module)
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

		ctrl.Get("", func(ctx core.Ctx) error {
			service := InjectModel[Book](module)
			data, err := service.Find(nil)
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
			Controllers: []core.Controller{bookController},
		})

		return bookMod
	}

	appModule := func() *core.DynamicModule {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Module{
				ForRoot(os.Getenv("MONGO_URI")),
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

	resp, err := testClient.Post(testServer.URL+"/app/books", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/app/books")
	require.Nil(t, err)
	require.Equal(t, 200, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/app/books/connect")
	require.Nil(t, err)
	require.Equal(t, 200, resp.StatusCode)
}
