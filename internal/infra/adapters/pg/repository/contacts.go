package repository

import (
	"github.com/AjxGnx/contacts-go/internal/domain/models"
	"gorm.io/gorm"
)

type Contacts interface {
	Create(contact models.Contact) (models.Contact, error)
}

type contacts struct {
	db *gorm.DB
}

func NewContacts(db *gorm.DB) Contacts {
	return &contacts{
		db,
	}
}

func (repo *contacts) Create(contact models.Contact) (models.Contact, error) {
	result := repo.db.Create(&contact).Scan(&contact)
	if result.Error != nil {
		return models.Contact{}, result.Error
	}

	return contact, nil
}
