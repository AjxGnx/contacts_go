package app

import (
	"github.com/AjxGnx/contacts-go/internal/domain/dto"
	"github.com/AjxGnx/contacts-go/internal/domain/models"
	"github.com/AjxGnx/contacts-go/internal/infra/adapters/pg/repository"
)

type Contacts interface {
	Create(contact dto.Contact) (models.Contact, error)
	GetByID(id uint) (models.Contact, error)
	Update(id uint, contact dto.Contact) (models.Contact, error)
	Delete(id uint) error
	Get(paginate dto.Paginate) (*models.Paginator, error)
}

type contacts struct {
	repo repository.Contacts
}

func NewContacts(repo repository.Contacts) Contacts {
	return &contacts{
		repo,
	}
}

func (app *contacts) Create(contact dto.Contact) (models.Contact, error) {
	return app.repo.Create(contact.ToModel())
}

func (app *contacts) GetByID(id uint) (models.Contact, error) {
	return app.repo.GetByID(id)
}

func (app *contacts) Update(id uint, contact dto.Contact) (models.Contact, error) {
	if _, err := app.GetByID(id); err != nil {
		return models.Contact{}, err
	}

	return app.repo.Update(id, contact.ToModel())
}

func (app *contacts) Delete(id uint) error {
	if _, err := app.GetByID(id); err != nil {
		return err
	}

	return app.repo.Delete(id)
}

func (app *contacts) Get(paginate dto.Paginate) (*models.Paginator, error) {
	return app.repo.Get(models.Paginator{Page: paginate.Page, Limit: paginate.Limit})
}
