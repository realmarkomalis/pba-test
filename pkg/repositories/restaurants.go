package repositories

import (
	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type IRestaurantsRepository interface {
	GetRestaurants() ([]*entities.Restaurant, error)
}

type RestaurantsRepository struct {
	DB *gorm.DB
}

func (r RestaurantsRepository) GetRestaurants() ([]*entities.Restaurant, error) {
	rs := []models.Restaurant{}

	err := r.DB.
		Model(&rs).
		Find(&rs).
		Error
	if err != nil {
		return nil, err
	}

	restaurants := []*entities.Restaurant{}
	for _, rest := range rs {
		r := rest.ModelToEntity()
		restaurants = append(restaurants, &r)
	}

	return restaurants, nil
}

func (r RestaurantsRepository) GetRestaurant(userID uint) (*entities.Restaurant, error) {
	rs := models.Restaurant{}

	err := r.DB.
		Model(&rs).
		Where("user_id = ?", userID).
		First(&rs).
		Error
	if err != nil {
		return nil, err
	}

	restaurant := rs.ModelToEntity()
	return &restaurant, nil
}

func (r RestaurantsRepository) GetUserRestaurants(userID uint) ([]entities.Restaurant, error) {
	rs := []models.Restaurant{}

	err := r.DB.
		Model(&rs).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("DropOffPoint").
		Find(&rs).
		Error
	if err != nil {
		return nil, err
	}

	restaurants := []entities.Restaurant{}
	for _, rest := range rs {
		r := rest.ModelToEntity()
		restaurants = append(restaurants, r)
	}

	return restaurants, nil
}
