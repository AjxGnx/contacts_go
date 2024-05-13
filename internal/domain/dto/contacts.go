package dto

import (
	"github.com/AjxGnx/contacts-go/internal/domain/models"
	"github.com/go-playground/validator/v10"
)

type Contact struct {
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

func (dto Contact) ToModel() models.Contact {
	return models.Contact{
		Name:        dto.Name,
		PhoneNumber: dto.PhoneNumber,
	}
}

func (dto Contact) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		return err
	}

	return nil
}
