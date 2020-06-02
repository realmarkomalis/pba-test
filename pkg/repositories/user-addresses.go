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
	return nil, nil
}
