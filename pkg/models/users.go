package models

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type User struct {
	gorm.Model
	FirstName     string
	LastName      string
	Email         string `gorm:"type:varchar(100);unique; not null"`
	IBAN          string
	PhoneNumber   string
	UserRole      UserRole
	UserRoleID    int
	UserAddresses []UserAddress
}

func (u User) ModelToEntity() entities.User {
	uas := []entities.UserAddress{}

	for _, ua := range u.UserAddresses {
		uas = append(uas, ua.ModelToEntity())
	}

	return entities.User{
		ID:            u.ID,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		UserAddresses: uas,
		UserRole:      u.UserRole.ModelToEntity(),
	}
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

func (u UserAddress) ModelToEntity() entities.UserAddress {
	return entities.UserAddress{
		ID:                u.ID,
		PostalCode:        u.PostalCode,
		StreetName:        u.StreetName,
		HouseNumber:       u.HouseNumber,
		HouseNumberSuffix: u.HouseNumberSuffix,
		City:              u.City,
		Country:           u.Country,
	}
}

type UserRole struct {
	gorm.Model
	Name     string `gorm:"not null"`
	SystemID string `gorm:"not null"`
}

func (u UserRole) ModelToEntity() entities.UserRole {
	return entities.UserRole{
		Name:     u.Name,
		SystemID: u.SystemID,
	}
}
