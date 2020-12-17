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

func (r RestaurantsRepository) GetRestaurants() ([]entities.Restaurant, error) {
	rs := []models.Restaurant{}

	err := r.DB.
		Model(&rs).
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

func (r RestaurantsRepository) GetRestaurant(rID uint) (*entities.Restaurant, error) {
	rs := models.Restaurant{}

	err := r.DB.
		Model(&rs).
		Where("id = ?", rID).
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

func (r RestaurantsRepository) GetSupplierRestaurants(userID uint) ([]entities.Restaurant, error) {
	ps := []models.PackageSupply{}

	err := r.DB.
		Model(&ps).
		Where("user_id = ?", userID).
		Preload("Restaurant").
		Find(&ps).
		Error
	if err != nil {
		return nil, err
	}

	seenRestaurant := make(map[uint]bool)
	restaurants := []entities.Restaurant{}
	for _, s := range ps {
		if _, ok := seenRestaurant[s.Restaurant.ID]; !ok {
			seenRestaurant[s.Restaurant.ID] = true
			restaurants = append(restaurants, s.Restaurant.ModelToEntity())
		}
	}

	return restaurants, nil
}
