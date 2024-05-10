package repository

import (
	"github.com/AjxGnx/contacts-go/internal/domain/models"
	"gorm.io/gorm"
)

type Contacts interface {
	Create(contact models.Contact) (models.Contact, error)
	GetByID(id uint) (models.Contact, error)
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

func (repo *contacts) GetByID(id uint) (models.Contact, error) {
	var contact models.Contact

	result := repo.db.First(&contact, id)
	if result.Error != nil {
		return contact, result.Error
	}

	return contact, nil
}
