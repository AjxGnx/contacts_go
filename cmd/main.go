package main

import (
	"fmt"
	"log"

	"github.com/AjxGnx/contacts-go/cmd/providers"
	"github.com/AjxGnx/contacts-go/config"
	"github.com/AjxGnx/contacts-go/internal/infra/api/router"
	"github.com/labstack/echo/v4"
)

func main() {
	container := providers.BuildContainer()
	err := container.Invoke(func(router *router.Router, server *echo.Echo) {
		router.Init()

		server.Logger.Fatal(server.Start(fmt.Sprintf("%s:%v", config.Environments().ServerHost,
			config.Environments().ServerPort)))
	})

	if err != nil {
		log.Panic(err)
	}
}
