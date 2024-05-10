package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AjxGnx/contacts-go/internal/app"
	"github.com/AjxGnx/contacts-go/internal/domain/dto"
	"github.com/labstack/echo/v4"
)

type Contacts interface {
	Create(ctx echo.Context) error
	GetByID(ctx echo.Context) error
	Update(ctx echo.Context) error
	Delete(ctx echo.Context) error
}

type contacts struct {
	app app.Contacts
}

func NewContacts(app app.Contacts) Contacts {
	return &contacts{
		app,
	}
}

func (handler *contacts) Create(ctx echo.Context) error {
	var contact dto.Contact

	if err := ctx.Bind(&contact); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result, err := handler.app.Create(contact)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("your contact number %s already exists",
				contact.PhoneNumber))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, dto.Message{
		Message: "contact created successfully",
		Data:    result,
	})
}

func (handler *contacts) GetByID(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	contact, err := handler.app.GetByID(uint(id))

	if err != nil {
		return errorValidator(err.Error(), id)
	}

	return ctx.JSON(http.StatusOK, dto.Message{
		Message: "contact successfully loaded",
		Data:    contact,
	})
}

func (handler *contacts) Update(ctx echo.Context) error {
	var contact dto.Contact

	contactID, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.Bind(&contact); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result, err := handler.app.Update(uint(contactID), contact)
	if err != nil {
		return errorValidator(err.Error(), contactID)
	}

	return ctx.JSON(http.StatusOK, dto.Message{
		Message: "contact updated successfully",
		Data:    result,
	})
}

func (handler *contacts) Delete(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := handler.app.Delete(uint(id)); err != nil {
		return errorValidator(err.Error(), id)
	}

	return ctx.JSON(http.StatusOK, dto.Message{
		Message: "contact successfully deleted",
	})
}

func errorValidator(errMessage string, id ...int) error {
	if errMessage == "record not found" {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("the contact: %v does not exist", id))
	}

	return echo.NewHTTPError(http.StatusInternalServerError, errMessage)
}
