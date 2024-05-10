package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AjxGnx/contacts-go/internal/domain/dto"
	"github.com/AjxGnx/contacts-go/internal/domain/models"
	mocks "github.com/AjxGnx/contacts-go/mocks/app"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type contactsTestSuite struct {
	suite.Suite
	Contacts  *mocks.Contacts
	underTest Contacts
}

func TestContactsSuite(t *testing.T) {
	suite.Run(t, new(contactsTestSuite))
}

func (suite *contactsTestSuite) SetupTest() {
	suite.Contacts = &mocks.Contacts{}
	suite.underTest = NewContacts(suite.Contacts)
}

func (suite *contactsTestSuite) TestCreate_WhenBindFail() {
	var httpError *echo.HTTPError

	body, _ := json.Marshal("")

	setupCase := SetupControllerCase(http.MethodPost, "/api/contacts/", bytes.NewBuffer(body))
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

	suite.Contacts.Mock.On("Create", contact).Return(models.Contact{}, err)

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

	suite.Contacts.Mock.On("Create", contact).Return(models.Contact{}, err)

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

	suite.Contacts.Mock.On("Create", contact).
		Return(models.Contact{Name: contact.Name, PhoneNumber: contact.PhoneNumber}, nil)

	suite.NoError(suite.underTest.Create(setupCase.context))
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
