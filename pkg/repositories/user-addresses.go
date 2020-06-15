package repositories

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type UserAddressesRepository struct {
	DB *gorm.DB
}

func (r UserAddressesRepository) CreateUserAddress(user *entities.User, address *entities.UserAddress) (*entities.UserAddress, error) {
	u := models.User{}
	err := r.DB.
		Where("id = ?", user.ID).
		First(&u).
		Error

	if err != nil {
		return nil, err
	}

	u.UserAddresses = []models.UserAddress{
		{
			UserID:            user.ID,
			PostalCode:        address.PostalCode,
			StreetName:        address.StreetName,
			HouseNumber:       address.HouseNumber,
			HouseNumberSuffix: address.HouseNumberSuffix,
			City:              address.City,
			Country:           address.Country,
		},
	}

	err = r.DB.Save(&u).Error
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (r UserAddressesRepository) UpdateUserAddress(user *entities.User, address *entities.UserAddress) (*entities.UserAddress, error) {
	a := models.UserAddress{}
	err := r.DB.
		FirstOrCreate(&a, models.UserAddress{
			UserID: user.ID,
		}).
		Error
	if err != nil {
		return nil, err
	}

	a.PostalCode = address.PostalCode
	a.StreetName = address.StreetName
	a.HouseNumber = address.HouseNumber
	a.HouseNumberSuffix = address.HouseNumberSuffix
	a.City = address.City
	a.Country = address.Country

	err = r.DB.
		Save(&a).
		Error
	if err != nil {
		return nil, err
	}

	return &entities.UserAddress{
		ID:                a.ID,
		PostalCode:        a.PostalCode,
		StreetName:        a.StreetName,
		HouseNumber:       a.HouseNumber,
		HouseNumberSuffix: a.HouseNumberSuffix,
		City:              a.City,
		Country:           a.Country,
	}, nil
}
