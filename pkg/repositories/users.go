package repositories

import (
	"errors"

	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type IUserRepository interface {
	Create(user *entities.User) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func (ur UserRepository) GetByEmail(email string) (*entities.User, error) {
	user := models.User{}
	role := models.UserRole{}
	addresses := []models.UserAddress{}
	err := ur.DB.
		First(&user, models.User{Email: email}).
		Related(&role).
		Related(&addresses).
		Error
	if err != nil {
		return nil, err
	}

	var as []entities.UserAddress
	if len(addresses) > 0 {
		as = append(as, entities.UserAddress{
			PostalCode:        addresses[0].PostalCode,
			StreetName:        addresses[0].StreetName,
			HouseNumber:       addresses[0].HouseNumber,
			HouseNumberSuffix: addresses[0].HouseNumberSuffix,
			City:              addresses[0].City,
			Country:           addresses[0].Country,
		})
	}

	return &entities.User{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		UserRole: entities.UserRole{
			Name:     role.Name,
			SystemID: role.SystemID,
		},
		UserAddresses: as,
	}, nil
}

func (ur UserRepository) Create(user *entities.User) (*entities.User, error) {
	u := models.User{}
	err := ur.DB.FirstOrCreate(&u, models.User{
		Email:      user.Email,
		UserRoleID: 4,
	}).Error
	if err != nil {
		return nil, err
	}

	if u.Email == "" {
		return nil, errors.New("Unable to create user in repository")
	}

	return &entities.User{
		ID:          u.ID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		PhoneNumber: u.PhoneNumber,
		UserRole: entities.UserRole{
			Name:     u.UserRole.Name,
			SystemID: u.UserRole.SystemID,
		},
	}, nil
}

func (ur UserRepository) UpdateUser(u *entities.User) (*entities.User, error) {
	user := models.User{}
	err := ur.DB.Model(&user).
		Where("id = ?", u.ID).
		Updates(models.User{
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
		}).
		Error
	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		UserRole: entities.UserRole{
			Name:     user.UserRole.Name,
			SystemID: user.UserRole.SystemID,
		},
	}, nil

}
