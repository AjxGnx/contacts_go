package dto

import (
	"github.com/AjxGnx/contacts-go/internal/domain/models"
)

type Contact struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

func (dto Contact) ToModel() models.Contact {
	return models.Contact{
		Name:        dto.Name,
		PhoneNumber: dto.PhoneNumber,
	}
}
