package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/AjxGnx/contacts-go/internal/domain/dto"
	"github.com/AjxGnx/contacts-go/internal/domain/models"
	mocks "github.com/AjxGnx/contacts-go/mocks/app"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type contactsTestSuite struct {
	suite.Suite
	app       *mocks.Contacts
	underTest Contacts
}

func TestContactsSuite(t *testing.T) {
	suite.Run(t, new(contactsTestSuite))
}

func (suite *contactsTestSuite) SetupTest() {
	suite.app = &mocks.Contacts{}
	suite.underTest = NewContacts(suite.app)
}

func (suite *contactsTestSuite) TestCreate_WhenBindFail() {
	var httpError *echo.HTTPError

	body, _ := json.Marshal("")

	setupCase := SetupControllerCase(http.MethodPost, "/api/contacts/", bytes.NewBuffer(body))
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.ErrorAs(suite.underTest.Create(setupCase.context), &httpError)
	suite.Equal(http.StatusBadRequest, httpError.Code)
}

func (suite *contactsTestSuite) TestCreate_WhenValidateFail() {
	var httpError *echo.HTTPError

	body := `{
		"name": "some string"
	}`

	setupCase := SetupControllerCase(http.MethodPost, "/api/exercise/numbers/", strings.NewReader(body))
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.ErrorAs(suite.underTest.Create(setupCase.context), &httpError)
	suite.Equal(http.StatusBadRequest, httpError.Code)
}

func (suite *contactsTestSuite) TestCreate_WhenFailByDuplicateContact() {
	var httpError *echo.HTTPError
	err := errors.New("SQLSTATE 23505")
	contact := dto.Contact{
		Name:        "test1",
		PhoneNumber: "+570000000",
	}
	body, _ := json.Marshal(contact)

	setupCase := SetupControllerCase(http.MethodPost, "/api/contacts/", bytes.NewBuffer(body))
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.app.Mock.On("Create", contact).Return(models.Contact{}, err)

	suite.ErrorAs(suite.underTest.Create(setupCase.context), &httpError)
	suite.Equal(http.StatusBadRequest, httpError.Code)
}

func (suite *contactsTestSuite) TestCreate_WhenFailByInternalError() {
	var httpError *echo.HTTPError
	err := errors.New("some error")
	contact := dto.Contact{
		Name:        "test2",
		PhoneNumber: "+570000001",
	}
	body, _ := json.Marshal(contact)

	setupCase := SetupControllerCase(http.MethodPost, "/api/contacts/", bytes.NewBuffer(body))
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.app.Mock.On("Create", contact).Return(models.Contact{}, err)

	suite.ErrorAs(suite.underTest.Create(setupCase.context), &httpError)
	suite.Equal(http.StatusInternalServerError, httpError.Code)
}

func (suite *contactsTestSuite) TestCreate_WhenSuccess() {
	contact := dto.Contact{
		Name:        "test3",
		PhoneNumber: "+570000002",
	}
	body, _ := json.Marshal(contact)

	setupCase := SetupControllerCase(http.MethodPost, "/api/contacts/", bytes.NewBuffer(body))
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.app.Mock.On("Create", contact).
		Return(models.Contact{Name: contact.Name, PhoneNumber: contact.PhoneNumber}, nil)

	suite.NoError(suite.underTest.Create(setupCase.context))
}

func (suite *contactsTestSuite) TestGetByID_WhenSuccess() {
	paramValue := 10
	param := "id"

	suite.app.Mock.On("GetByID", uint(paramValue)).
		Return(models.Contact{ID: 10}, nil)

	setupCase := SetupControllerCase(http.MethodPost, "/api/contacts/10", nil)
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	setupCase.context.SetParamNames(param)
	setupCase.context.SetParamValues(strconv.Itoa(paramValue))

	suite.NoError(suite.underTest.GetByID(setupCase.context))
	suite.Equal(http.StatusOK, setupCase.Res.Code)
}

func (suite *contactsTestSuite) TestGetByID_WhenContactNotFound() {
	var httpError *echo.HTTPError

	paramValue := 10
	param := "id"
	expectedError := errors.New("record not found")

	suite.app.Mock.On("GetByID", uint(paramValue)).
		Return(models.Contact{}, expectedError)

	setupCase := SetupControllerCase(http.MethodPost, "/api/contacts/10", nil)
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	setupCase.context.SetParamNames(param)
	setupCase.context.SetParamValues(strconv.Itoa(paramValue))

	suite.ErrorAs(suite.underTest.GetByID(setupCase.context), &httpError)
	suite.Equal(http.StatusNotFound, httpError.Code)
}

func (suite *contactsTestSuite) TestGetByID_WhenFail() {
	var httpError *echo.HTTPError

	paramValue := 10
	param := "id"
	expectedError := errors.New("some error")

	suite.app.Mock.On("GetByID", uint(paramValue)).
		Return(models.Contact{}, expectedError)

	setupCase := SetupControllerCase(http.MethodPost, "/api/contacts/10", nil)
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	setupCase.context.SetParamNames(param)
	setupCase.context.SetParamValues(strconv.Itoa(paramValue))

	suite.ErrorAs(suite.underTest.GetByID(setupCase.context), &httpError)
	suite.Equal(http.StatusInternalServerError, httpError.Code)
}

func (suite *contactsTestSuite) TestUpdate_WhenSuccess() {
	paramValue := 10
	param := "id"

	contact := dto.Contact{
		Name:        "test3",
		PhoneNumber: "+570000002",
	}

	body, _ := json.Marshal(contact)

	suite.app.Mock.On("Update", uint(paramValue), contact).
		Return(models.Contact{ID: 10}, nil)

	setupCase := SetupControllerCase(http.MethodPut, "/api/contacts/10", bytes.NewBuffer(body))
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	setupCase.context.SetParamNames(param)
	setupCase.context.SetParamValues(strconv.Itoa(paramValue))

	suite.NoError(suite.underTest.Update(setupCase.context))
	suite.Equal(http.StatusOK, setupCase.Res.Code)
}

func (suite *contactsTestSuite) TestUpdate_WhenContactNotFound() {
	var httpError *echo.HTTPError

	paramValue := 10
	param := "id"

	contact := dto.Contact{
		Name:        "test3",
		PhoneNumber: "+570000002",
	}

	expectedError := errors.New("record not found")

	body, _ := json.Marshal(contact)

	suite.app.Mock.On("Update", uint(paramValue), contact).
		Return(models.Contact{}, expectedError)

	setupCase := SetupControllerCase(http.MethodPut, "/api/contacts/10", bytes.NewBuffer(body))
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	setupCase.context.SetParamNames(param)
	setupCase.context.SetParamValues(strconv.Itoa(paramValue))

	suite.ErrorAs(suite.underTest.Update(setupCase.context), &httpError)
	suite.Equal(http.StatusNotFound, httpError.Code)
}

func (suite *contactsTestSuite) TestUpdate_WhenFail() {
	var httpError *echo.HTTPError

	paramValue := 10
	param := "id"

	contact := dto.Contact{
		Name:        "test3",
		PhoneNumber: "+570000002",
	}

	expectedError := errors.New("some error")

	body, _ := json.Marshal(contact)

	suite.app.Mock.On("Update", uint(paramValue), contact).
		Return(models.Contact{}, expectedError)

	setupCase := SetupControllerCase(http.MethodPut, "/api/contacts/10", bytes.NewBuffer(body))
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	setupCase.context.SetParamNames(param)
	setupCase.context.SetParamValues(strconv.Itoa(paramValue))

	suite.ErrorAs(suite.underTest.Update(setupCase.context), &httpError)
	suite.Equal(http.StatusInternalServerError, httpError.Code)
}

func (suite *contactsTestSuite) TestUpdate_WhenBindFail() {
	var httpError *echo.HTTPError

	body, _ := json.Marshal("")

	setupCase := SetupControllerCase(http.MethodPost, "/api/contacts/", bytes.NewBuffer(body))
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.ErrorAs(suite.underTest.Update(setupCase.context), &httpError)
	suite.Equal(http.StatusBadRequest, httpError.Code)
}
func (suite *contactsTestSuite) TestUpdate_WhenValidateFail() {
	var httpError *echo.HTTPError

	paramValue := 10
	param := "id"

	contact := dto.Contact{
		Name:        "test3",
		PhoneNumber: "",
	}

	body, _ := json.Marshal(contact)

	setupCase := SetupControllerCase(http.MethodPut, "/api/contacts/10", bytes.NewBuffer(body))
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	setupCase.context.SetParamNames(param)
	setupCase.context.SetParamValues(strconv.Itoa(paramValue))

	suite.ErrorAs(suite.underTest.Update(setupCase.context), &httpError)
	suite.Equal(http.StatusBadRequest, httpError.Code)
}

func (suite *contactsTestSuite) TestDelete_WhenSuccess() {
	paramValue := 10
	param := "id"

	suite.app.Mock.On("Delete", uint(paramValue)).
		Return(nil)

	setupCase := SetupControllerCase(http.MethodDelete, "/api/contacts/10", nil)
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	setupCase.context.SetParamNames(param)
	setupCase.context.SetParamValues(strconv.Itoa(paramValue))

	suite.NoError(suite.underTest.Delete(setupCase.context))
	suite.Equal(http.StatusOK, setupCase.Res.Code)
}

func (suite *contactsTestSuite) TestDelete_WhenContactNotFound() {
	var httpError *echo.HTTPError

	paramValue := 10
	param := "id"
	expectedError := errors.New("record not found")

	suite.app.Mock.On("Delete", uint(paramValue)).
		Return(expectedError)

	setupCase := SetupControllerCase(http.MethodDelete, "/api/contacts/10", nil)
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	setupCase.context.SetParamNames(param)
	setupCase.context.SetParamValues(strconv.Itoa(paramValue))

	suite.ErrorAs(suite.underTest.Delete(setupCase.context), &httpError)
	suite.Equal(http.StatusNotFound, httpError.Code)
}

func (suite *contactsTestSuite) TestDelete_WhenFail() {
	var httpError *echo.HTTPError

	paramValue := 10
	param := "id"
	expectedError := errors.New("some error")

	suite.app.Mock.On("Delete", uint(paramValue)).
		Return(expectedError)

	setupCase := SetupControllerCase(http.MethodDelete, "/api/contacts/10", nil)
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	setupCase.context.SetParamNames(param)
	setupCase.context.SetParamValues(strconv.Itoa(paramValue))

	suite.ErrorAs(suite.underTest.Delete(setupCase.context), &httpError)
	suite.Equal(http.StatusInternalServerError, httpError.Code)
}

func (suite *contactsTestSuite) TestGet_WhenSuccess() {
	paginateValues := dto.Paginate{
		Page:  1,
		Limit: 10,
	}

	suite.app.Mock.On("Get", paginateValues).
		Return(&models.Paginator{}, nil)

	setupCase := SetupControllerCase(http.MethodGet, "/api/contacts/?page=1&limit=10", nil)
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.NoError(suite.underTest.Get(setupCase.context))
	suite.Equal(http.StatusOK, setupCase.Res.Code)
}

func (suite *contactsTestSuite) TestGet_WhenFail() {
	var httpError *echo.HTTPError

	paginateValues := dto.Paginate{
		Page:  1,
		Limit: 10,
	}
	expectedError := errors.New("some error")

	suite.app.Mock.On("Get", paginateValues).
		Return(&models.Paginator{}, expectedError)

	setupCase := SetupControllerCase(http.MethodGet, "/api/contacts/?page=1&limit=10", nil)
	setupCase.Req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.ErrorAs(suite.underTest.Get(setupCase.context), &httpError)
	suite.Equal(http.StatusInternalServerError, httpError.Code)
}

type ControllerCase struct {
	Req     *http.Request
	Res     *httptest.ResponseRecorder
	context echo.Context
}

func SetupControllerCase(method string, url string, body io.Reader) ControllerCase {
	engine := echo.New()
	req := httptest.NewRequest(method, url, body)
	res := httptest.NewRecorder()
	ctxEngine := engine.NewContext(req, res)

	return ControllerCase{req, res, ctxEngine}
}
