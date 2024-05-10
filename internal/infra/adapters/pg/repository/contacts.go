package repository

import (
	"math"

	"github.com/AjxGnx/contacts-go/internal/domain/models"
	"gorm.io/gorm"
)

type Contacts interface {
	Create(contact models.Contact) (models.Contact, error)
	GetByID(id uint) (models.Contact, error)
	Update(id uint, account models.Contact) (models.Contact, error)
	Delete(id uint) error
	Get(paginate models.Paginator) (*models.Paginator, error)
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

func (repo *contacts) Update(id uint, contact models.Contact) (models.Contact, error) {
	result := repo.db.
		Model(&contact).
		Where("id = ?", id).
		Updates(contact).
		Scan(&contact)

	if result.Error != nil {
		return contact, result.Error
	}

	return contact, nil
}

func (repo *contacts) Delete(id uint) error {
	result := repo.db.
		Where("id = ?", id).
		Delete(&models.Contact{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *contacts) Get(paginate models.Paginator) (*models.Paginator, error) {
	var contacts []models.Contact

	offset := (paginate.Page - 1) * paginate.Limit

	err := repo.db.Offset(offset).Limit(paginate.Limit).Find(&contacts).Error
	if err != nil {
		return nil, err
	}

	totalRecords, err := repo.countTotalRecords()
	if err != nil {
		return nil, err
	}

	paginator := &models.Paginator{
		TotalRecord: totalRecords,
		TotalPage:   int(math.Ceil(float64(totalRecords) / float64(paginate.Limit))),
		Records:     contacts,
		Offset:      offset,
		Limit:       paginate.Limit,
		Page:        paginate.Page,
	}

	if paginate.Page > 1 {
		paginator.PrevPage = paginate.Page - 1
	} else {
		paginator.PrevPage = paginate.Page
	}

	paginator.NextPage = paginate.Page + 1

	return paginator, nil

}

func (repo *contacts) countTotalRecords() (int64, error) {
	var total int64

	if err := repo.db.Model(&models.Contact{}).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
