package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	FirstName     string
	LastName      string
	Email         string `gorm:"type:varchar(100);unique; not null"`
	IBAN          string
	Role          UserRole
	UserRoleID    int
	UserAddresses []UserAddress
}

type UserAddress struct {
	gorm.Model
	UserID            uint
	PostalCode        string `gorm:"not null"`
	StreetName        string `gorm:"not null"`
	HouseNumber       string `gorm:"not null"`
	HouseNumberSuffix string
	City              string `gorm:"not null"`
	Country           string `gorm:"not null"`
}

type UserRole struct {
	gorm.Model
	Name     string `gorm:"not null"`
	SystemID string `gorm:"not null"`
}
