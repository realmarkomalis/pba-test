package usecases

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type IReturnsUsecase interface {
	CreateNewReturn(userID, packageID uint) (*entities.Return, error)
	CreateReturnRequest(userID, packageID uint) (*entities.Return, error)
	CreatePackagePickup(userID, packageID uint) (*entities.PackagePickup, error)
}

type ReturnsUsecase struct {
	ReturnsRepo repositories.IReturnRepository
}

func (u ReturnsUsecase) CreateNewReturn(userID, packageID uint) (*entities.Return, error) {
	ret, err := u.ReturnsRepo.Create(userID, packageID)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (u ReturnsUsecase) CreateReturnRequest(userID, packageID uint) (*entities.Return, error) {
	return nil, nil
}

func (u ReturnsUsecase) CreatePackagePickup(userID, packageID uint) (*entities.PackagePickup, error) {
	pickUp, err := u.ReturnsRepo.CreatePackagePickup(userID, packageID)
	if err != nil {
		return nil, err
	}

	return pickUp, nil
}
