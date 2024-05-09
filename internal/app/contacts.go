package app

import (
	"github.com/AjxGnx/contacts-go/internal/domain/dto"
	"github.com/AjxGnx/contacts-go/internal/domain/models"
	"github.com/AjxGnx/contacts-go/internal/infra/adapters/pg/repository"
)

type Contacts interface {
	Create(contact dto.Contact) (models.Contact, error)
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
