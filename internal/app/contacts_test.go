package app

import (
	"errors"
	"testing"

	"github.com/AjxGnx/contacts-go/internal/domain/dto"
	"github.com/AjxGnx/contacts-go/internal/domain/models"
	mocks "github.com/AjxGnx/contacts-go/mocks/infra/adapters/pg/repository"
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

func (suite *contactsTestSuite) TestCreate_WhenSuccess() {
	contact := dto.Contact{
		Name:        "test",
		PhoneNumber: "+570000000",
	}

	expected := models.Contact{Name: contact.Name, PhoneNumber: contact.PhoneNumber, ID: 1}

	suite.Contacts.Mock.On("Create", models.Contact{
		Name:        contact.Name,
		PhoneNumber: contact.PhoneNumber,
	}).Return(expected, nil)

	contactModel, err := suite.underTest.Create(contact)

	suite.NoError(err)
	suite.Equal(expected, contactModel)
}

func (suite *contactsTestSuite) TestCreate_WhenFail() {
	contact := dto.Contact{
		Name:        "test",
		PhoneNumber: "+570000000",
	}

	expectedError := errors.New("some error")

	suite.Contacts.Mock.On("Create", models.Contact{
		Name:        contact.Name,
		PhoneNumber: contact.PhoneNumber,
	}).Return(models.Contact{}, expectedError)

	contactModel, err := suite.underTest.Create(contact)

	suite.Error(err)
	suite.Equal(models.Contact{}, contactModel)
}

func (suite *contactsTestSuite) TestGetByID_WhenSuccess() {
	expected := models.Contact{Name: "test", PhoneNumber: "+570000000", ID: 1}

	suite.Contacts.Mock.On("GetByID", uint(1)).Return(expected, nil)

	contactModel, err := suite.underTest.GetByID(uint(1))

	suite.NoError(err)
	suite.Equal(expected, contactModel)
}

func (suite *contactsTestSuite) TestGetByID_WhenFail() {
	expectedError := errors.New("some error")

	suite.Contacts.Mock.On("GetByID", uint(1)).Return(models.Contact{}, expectedError)

	contactModel, err := suite.underTest.GetByID(uint(1))

	suite.Error(err)
	suite.Equal(models.Contact{}, contactModel)
}
