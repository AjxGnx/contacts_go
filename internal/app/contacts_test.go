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
	repo      *mocks.Contacts
	underTest Contacts
}

func TestContactsSuite(t *testing.T) {
	suite.Run(t, new(contactsTestSuite))
}

func (suite *contactsTestSuite) SetupTest() {
	suite.repo = &mocks.Contacts{}
	suite.underTest = NewContacts(suite.repo)
}

func (suite *contactsTestSuite) TestCreate_WhenSuccess() {
	contact := dto.Contact{
		Name:        "test",
		PhoneNumber: "+570000000",
	}

	expected := models.Contact{Name: contact.Name, PhoneNumber: contact.PhoneNumber, ID: 1}

	suite.repo.Mock.On("Create", models.Contact{
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

	suite.repo.Mock.On("Create", models.Contact{
		Name:        contact.Name,
		PhoneNumber: contact.PhoneNumber,
	}).Return(models.Contact{}, expectedError)

	contactModel, err := suite.underTest.Create(contact)

	suite.Error(err)
	suite.Equal(models.Contact{}, contactModel)
}

func (suite *contactsTestSuite) TestGetByID_WhenSuccess() {
	expected := models.Contact{Name: "test", PhoneNumber: "+570000000", ID: 1}

	suite.repo.Mock.On("GetByID", uint(1)).Return(expected, nil)

	contactModel, err := suite.underTest.GetByID(uint(1))

	suite.NoError(err)
	suite.Equal(expected, contactModel)
}

func (suite *contactsTestSuite) TestGetByID_WhenFail() {
	expectedError := errors.New("some error")

	suite.repo.Mock.On("GetByID", uint(1)).Return(models.Contact{}, expectedError)

	contactModel, err := suite.underTest.GetByID(uint(1))

	suite.Error(err)
	suite.Equal(models.Contact{}, contactModel)
}

func (suite *contactsTestSuite) TestUpdate_WhenSuccess() {
	contact := dto.Contact{
		Name:        "test",
		PhoneNumber: "+570000000",
	}

	expected := models.Contact{Name: contact.Name, PhoneNumber: contact.PhoneNumber, ID: 1}

	suite.repo.Mock.On("GetByID", uint(1)).Return(models.Contact{}, nil)
	suite.repo.Mock.On("Update", uint(1), models.Contact{
		Name:        contact.Name,
		PhoneNumber: contact.PhoneNumber,
	}).Return(expected, nil)

	contactModel, err := suite.underTest.Update(uint(1), contact)

	suite.NoError(err)
	suite.Equal(expected, contactModel)
}

func (suite *contactsTestSuite) TestUpdate_WhenFail() {
	contact := dto.Contact{
		Name:        "test",
		PhoneNumber: "+570000000",
	}

	expectedError := errors.New("some error")

	suite.repo.Mock.On("GetByID", uint(1)).Return(models.Contact{}, nil)
	suite.repo.Mock.On("Update", uint(1), models.Contact{
		Name:        contact.Name,
		PhoneNumber: contact.PhoneNumber,
	}).Return(models.Contact{}, expectedError)

	contactModel, err := suite.underTest.Update(uint(1), contact)

	suite.Error(err)
	suite.Equal(models.Contact{}, contactModel)
}

func (suite *contactsTestSuite) TestUpdate_WhenGetByIDFail() {
	contact := dto.Contact{
		Name:        "test",
		PhoneNumber: "+570000000",
	}
	expectedError := errors.New("some error")

	suite.repo.Mock.On("GetByID", uint(1)).Return(models.Contact{}, expectedError)

	contactModel, err := suite.underTest.Update(uint(1), contact)

	suite.Error(err)
	suite.Equal(models.Contact{}, contactModel)
}

func (suite *contactsTestSuite) TestDelete_WhenSuccess() {
	suite.repo.Mock.On("GetByID", uint(1)).Return(models.Contact{}, nil)
	suite.repo.Mock.On("Delete", uint(1)).Return(nil)

	suite.NoError(suite.underTest.Delete(uint(1)))
}

func (suite *contactsTestSuite) TestDelete_WhenGetByIDFail() {
	expectedError := errors.New("some error")

	suite.repo.Mock.On("GetByID", uint(1)).Return(models.Contact{}, expectedError)

	suite.Error(suite.underTest.Delete(uint(1)))
}

func (suite *contactsTestSuite) TestDelete_WhenFail() {
	expectedError := errors.New("some error")

	suite.repo.Mock.On("GetByID", uint(1)).Return(models.Contact{}, nil)
	suite.repo.Mock.On("Delete", uint(1)).Return(expectedError)

	suite.Error(suite.underTest.Delete(uint(1)))
}

func (suite *contactsTestSuite) TestGet_WhenSuccess() {
	paginate := dto.Paginate{
		Page:  1,
		Limit: 10,
	}
	suite.repo.Mock.On("Get", models.Paginator{Page: paginate.Page, Limit: paginate.Limit}).
		Return(&models.Paginator{}, nil)

	_, err := suite.underTest.Get(paginate)

	suite.NoError(err)
}

func (suite *contactsTestSuite) TestGet_WhenFail() {
	paginate := dto.Paginate{
		Page:  1,
		Limit: 10,
	}
	expectedError := errors.New("some error")
	suite.repo.Mock.On("Get", models.Paginator{Page: paginate.Page, Limit: paginate.Limit}).
		Return(&models.Paginator{}, expectedError)

	_, err := suite.underTest.Get(paginate)

	suite.Error(err)
}
