package router

import (
	_ "github.com/AjxGnx/contacts-go/docs"
	"github.com/AjxGnx/contacts-go/internal/infra/api/handler"
	"github.com/AjxGnx/contacts-go/internal/infra/api/router/group"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	server        *echo.Echo
	contactsGroup group.Contacts
}

func New(
	server *echo.Echo,
	contactsGroup group.Contacts,
) *Router {
	return &Router{
		server,
		contactsGroup,
	}
}

func (router *Router) Init() {
	router.server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} latency=${latency_human}\n",
	}))

	router.server.Use(middleware.Recover())

	basePath := router.server.Group("/api")

	basePath.GET("/swagger/*", echoSwagger.WrapHandler)
	basePath.GET("/health", handler.HealthCheck)

	router.contactsGroup.Resource(basePath)
}
