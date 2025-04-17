package mongoose

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Module(t *testing.T) {
	type Book struct {
		BaseSchema `bson:"inline"`
		Title      string `bson:"title"`
		Author     string `bson:"author"`
	}
	bookModel := NewModel[Book]("Book")

	bookController := func(module core.Module) core.Controller {
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

	bookModule := func(module core.Module) core.Module {
		bookMod := module.New(core.NewModuleOptions{
			Imports:     []core.Modules{ForFeature(bookModel)},
			Controllers: []core.Controllers{bookController},
		})

		return bookMod
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				ForRoot(os.Getenv("MONGO_URI") + "/test"),
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
