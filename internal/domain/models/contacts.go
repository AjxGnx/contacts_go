package models

type Contact struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"not null"`
	PhoneNumber string `json:"phone_number" gorm:"unique;not null"`
}
