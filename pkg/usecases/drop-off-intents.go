package usecases

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type DropOffIntentUsecases struct {
	ReturnsRepo        repositories.IReturnRepository
	DropOffIntentsRepo repositories.DropOffIntentsRepository
}

func (u DropOffIntentUsecases) CreateDropOffIntent(userID, packageID, dropOffPointID uint) (*entities.DropOffIntent, error) {
	rs, err := u.ReturnsRepo.GetActiveReturns(packageID)

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
	if r.Status != models.Created.String() && r.Status != models.Dispatched.String() {
		return nil, entities.APIError{
			Message: "Unable to create a drop-off intent for this package_id",
			Code:    "0",
		}
	}

	di, err := u.DropOffIntentsRepo.CreateDropOffIntent(userID, r.ID, dropOffPointID)
	if err != nil {
		return nil, entities.APIError{

			Message: err.Error(),
			Code:    "0",
		}
	}

	_, err = u.ReturnsRepo.SetReturnStatus(r.ID, models.DropOffIntended)
	if err != nil {
		return nil, entities.APIError{
			Message: err.Error(),
			Code:    "0",
		}
	}

	return di, nil
}
