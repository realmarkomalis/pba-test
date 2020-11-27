package usecases

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type DropOffCollectUsecases struct {
	ReturnsRepository         repositories.IReturnRepository
	RestaurantsRepository     repositories.RestaurantsRepository
	DropOffCollectsRepository repositories.DropOffCollectsRepository
}

func (u DropOffCollectUsecases) CreateDropOffCollect(userID, packageID, dropOffPointID uint) (*entities.DropOffCollect, error) {
	rs, err := u.ReturnsRepository.GetActiveReturns(packageID)

	if len(rs) == 0 {
		return nil, entities.APIError{
			Message: "No active return for this package_id",
			Code:    "0",
		}
	}

	if len(rs) > 1 {
		return nil, entities.APIError{
			Message: "Inconsistent state found for this package_id",
			Code:    "0",
		}
	}

	r := rs[0]
	if !(r.Status == models.Created.String() || r.Status == models.Dispatched.String() || r.Status == models.DropOffIntended.String()) {
		return nil, entities.APIError{
			Message: "Unable to create a drop-off collect for this package_id",
			Code:    "0",
		}
	}

	dc, err := u.DropOffCollectsRepository.CreateDropOffCollect(userID, r.ID, dropOffPointID)
	if err != nil {
		return nil, entities.APIError{

			Message: err.Error(),
			Code:    "0",
		}
	}

	_, err = u.ReturnsRepository.SetReturnStatus(r.ID, models.Collected)
	if err != nil {
		return nil, entities.APIError{
			Message: err.Error(),
			Code:    "0",
		}
	}

	return dc, nil
}

func (u DropOffCollectUsecases) CreateDropOffCollectAndSupply(userID, packageID, dropOffPointID uint) (*entities.DropOffCollect, error) {
	dc, err := u.CreateDropOffCollect(userID, packageID, dropOffPointID)
	if err != nil {
		return nil, entities.APIError{
			Message: err.Error(),
			Code:    "0",
		}
	}

	r, err := u.ReturnsRepository.CreateReturn(packageID)
	if err != nil {
		return nil, entities.APIError{
			Message: err.Error(),
			Code:    "0",
		}
	}

	resto, err := u.RestaurantsRepository.GetRestaurant(userID)
	if err != nil {
		return nil, entities.APIError{
			Message: err.Error(),
			Code:    "0",
		}
	}

	r, err = u.ReturnsRepository.CreatePackageSupply(r.ID, userID, resto.ID)
	if err != nil {
		return nil, entities.APIError{
			Message: err.Error(),
			Code:    "0",
		}
	}

	return dc, nil
}
