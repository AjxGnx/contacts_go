package providers

import (
	"github.com/AjxGnx/contacts-go/internal/infra/api/handler"
	"github.com/AjxGnx/contacts-go/internal/infra/api/router"
	"github.com/AjxGnx/contacts-go/internal/infra/api/router/group"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

var Container *dig.Container

func BuildContainer() *dig.Container {
	Container = dig.New()

	_ = Container.Provide(func() *echo.Echo {
		return echo.New()
	})

	_ = Container.Provide(handler.NewContactsHandler)
	_ = Container.Provide(group.NewContactsGroup)
	_ = Container.Provide(router.New)

	return Container
}
