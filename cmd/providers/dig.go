package providers

import (
	"github.com/AjxGnx/contacts-go/internal/app"
	"github.com/AjxGnx/contacts-go/internal/infra/adapters/pg"
	"github.com/AjxGnx/contacts-go/internal/infra/adapters/pg/repository"
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

	_ = Container.Provide(router.New)
	_ = Container.Provide(pg.ConnInstance)

	_ = Container.Provide(group.NewContacts)
	_ = Container.Provide(handler.NewContacts)
	_ = Container.Provide(app.NewContacts)
	_ = Container.Provide(repository.NewContacts)

	return Container
}
