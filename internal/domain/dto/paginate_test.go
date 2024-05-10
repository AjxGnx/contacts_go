package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginate_SetDefaultLimitAndPage(t *testing.T) {
	paginate := Paginate{}
	paginate.SetDefaultLimitAndPage()

	expectedPaginate := Paginate{
		Page:  1,
		Limit: 10,
	}

	assert.Equal(t, expectedPaginate, paginate)
}
