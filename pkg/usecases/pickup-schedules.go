package usecases

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type IPickupScheduleUsecase interface {
	GetScheduledReturnEntries() ([]*entities.UserReturnEntry, error)
}

type PickupScheduleUsecase struct {
	// ReturnsRepo     repositories.IReturnRepository
	UserReturnRepo repositories.IUserReturnRepository
}

func (u PickupScheduleUsecase) GetScheduledReturnEntries() ([]*entities.UserReturnEntry, error) {
	rets, err := u.UserReturnRepo.GetScheduledReturnEntries()
	if err != nil {
		return nil, err
	}

	e := []*entities.UserReturnEntry{}
	for _, ee := range rets {
		if ee.User.UserRole.SystemID == "customer" {
			e = append(e, ee)
		}
	}

	return e, nil
}
