package usecases

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type ReportsUsecase struct {
	ReportsRepository     repositories.ReportsRepository
	RestaurantsRepository repositories.RestaurantsRepository
}

func (r ReportsUsecase) GetRepotableRestaurants(userID uint) ([]entities.Restaurant, error) {
	rs, err := r.RestaurantsRepository.GetSupplierRestaurants(userID)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (r ReportsUsecase) GetUsersRestaurants(userID uint) ([]entities.Restaurant, error) {
	rs, err := r.RestaurantsRepository.GetUserRestaurants(userID)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (r ReportsUsecase) GetSupplierTotalsReport(userID uint, startInterval, endInterval string) (*entities.SupplierReport, error) {
	rs, err := r.RestaurantsRepository.GetSupplierRestaurants(userID)
	if err != nil {
		return nil, err
	}

	rp, err := r.ReportsRepository.SupplierTotalsReport(rs, userID, startInterval, endInterval)
	if err != nil {
		return nil, err
	}

	return rp, nil
}

func (r ReportsUsecase) GetSupplierRestaurantReport(userID, restaurantID uint, startInterval, endInterval string) (*entities.RestaurantReport, error) {
	rs, err := r.RestaurantsRepository.GetRestaurant(restaurantID)
	if err != nil {
		return nil, err
	}

	rp, err := r.ReportsRepository.SupplierRestaurantReport(*rs, userID, startInterval, endInterval)
	if err != nil {
		return nil, err
	}

	return rp, nil
}

func (r ReportsUsecase) GetRestaurantReport(userID, restaurantID uint, startInterval, endInterval string) (*entities.RestaurantReport, error) {
	rs, err := r.RestaurantsRepository.GetRestaurant(restaurantID)
	if err != nil {
		return nil, err
	}

	rp, err := r.ReportsRepository.RestaurantReport(*rs, startInterval, endInterval)
	if err != nil {
		return nil, err
	}

	return rp, nil
}
