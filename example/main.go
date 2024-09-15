package main

import (
	"github.com/tinh-tinh/mongoose/example/app"
	"github.com/tinh-tinh/tinhtinh/core"
)

func main() {
	server := core.CreateFactory(app.NewModule, "api")

	server.Listen(3000)
}
