package dto

import (
	"testing"

	"github.com/AjxGnx/contacts-go/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestContact_ToModel(t *testing.T) {
	contact := Contact{
		Name:        "test",
		PhoneNumber: "+570000000",
	}

	contactExpected := models.Contact{
		Name:        "test",
		PhoneNumber: "+570000000",
	}

	assert.Equal(t, contactExpected, contact.ToModel())
}

func TestContact_Validate(t *testing.T) {
	assert.Error(t, Contact{}.Validate())

	assert.NoError(t, Contact{
		Name:        "name",
		PhoneNumber: "phone number",
	}.Validate())
}
