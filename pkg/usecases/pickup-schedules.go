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
	return u.UserReturnRepo.GetScheduledReturnEntries()
}
