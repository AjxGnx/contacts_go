package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Contacts interface {
	Create(ctx echo.Context) error
}

type contacts struct {
}

func NewContactsHandler() Contacts {
	return contacts{}
}

func (c contacts) Create(ctx echo.Context) error {
	fmt.Println("metodo para crear un contacto")

	return nil
}
