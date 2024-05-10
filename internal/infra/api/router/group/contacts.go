package group

import (
	"github.com/AjxGnx/contacts-go/internal/infra/api/handler"
	"github.com/labstack/echo/v4"
)

const contactsPath = "/contacts/"

type Contacts interface {
	Resource(c *echo.Group)
}

type contacts struct {
	handler handler.Contacts
}

func NewContacts(handler handler.Contacts) Contacts {
	return &contacts{
		handler,
	}
}

func (routes *contacts) Resource(c *echo.Group) {
	groupPath := c.Group(contactsPath)
	groupPath.POST("", routes.handler.Create)
	groupPath.GET(":id", routes.handler.GetByID)
	groupPath.PUT(":id", routes.handler.Update)
}
