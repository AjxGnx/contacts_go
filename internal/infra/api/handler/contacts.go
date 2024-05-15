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
	Get(ctx echo.Context) error
}

type contacts struct {
	app app.Contacts
}

func NewContacts(app app.Contacts) Contacts {
	return &contacts{
		app,
	}
}

// @Tags         Contacts
// @Summary      Create a contact
// @Description  Create a contact
// @Accept       json
// @Produce      json
// @Param        request  body      dto.Contact  true  "Request Body"
// @Success      200      {object}  dto.Message{data=models.Contact}
// @Failure      400      {object}  dto.MessageError
// @Failure      500      {object}  dto.MessageError
// @Router       /contacts/ [post]
func (handler *contacts) Create(ctx echo.Context) error {
	var contact dto.Contact

	if err := ctx.Bind(&contact); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := contact.Validate(); err != nil {
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

// @Tags         Contacts
// @Summary      Get Contact by id
// @Description  Get Contact by id
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "value of record to find"
// @Success      200      {object}  models.Contact
// @Failure      404      {object}  dto.MessageError
// @Failure      500      {object}  dto.MessageError
// @Router       /contacts/{id} [get]
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

// @Tags         Contacts
// @Summary      Update Contact by id
// @Description  Update Contact by id
// @Accept       json
// @Produce      json
// @Param        request  body      dto.Contact  true  "Request Body"
// @Param        id       path      int          true  "value of record to update"
// @Success      200  {object}  models.Contact
// @Failure      404  {object}  dto.MessageError
// @Failure      500  {object}  dto.MessageError
// @Router       /contacts/{id} [put]
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

// @Tags         Contacts
// @Summary      Delete Contact by id
// @Description  Delete Contact by id
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "value of record to delete"
// @Success      200  {object}  dto.Message{}
// @Failure      404  {object}  dto.MessageError
// @Failure      500  {object}  dto.MessageError
// @Router       /contacts/{id} [delete]
func (handler *contacts) Delete(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := handler.app.Delete(uint(id)); err != nil {
		return errorValidator(err.Error(), id)
	}

	return ctx.JSON(http.StatusOK, dto.Message{
		Message: "contact successfully deleted",
	})
}

// @Tags         Contacts
// @Summary      Get contacts
// @Description  Get contacts using pagination
// @Accept       json
// @Produce      json
// @Param        limit  query     string  true  "limit to find records"
// @Param        page   query     string  true  "page to find records"
// @Success      200    {object}  dto.Message{data=models.Paginator{records=[]models.Contact}}
// @Failure      500    {object}  dto.MessageError
// @Router       /contacts/ [get]
func (handler *contacts) Get(context echo.Context) error {
	page, _ := strconv.Atoi(context.QueryParam("page"))
	limit, _ := strconv.Atoi(context.QueryParam("limit"))
	paginate := dto.Paginate{
		Page:  page,
		Limit: limit,
	}
	paginate.SetDefaultLimitAndPage()

	categorizations, err := handler.app.Get(paginate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, dto.Message{
		Message: fmt.Sprintf("contacts successfully loaded"),
		Data:    categorizations,
	})

}

func errorValidator(errMessage string, id ...int) error {
	if errMessage == "record not found" {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("the contact: %v does not exist", id))
	}

	return echo.NewHTTPError(http.StatusInternalServerError, errMessage)
}
