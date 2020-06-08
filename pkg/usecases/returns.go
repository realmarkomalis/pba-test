package usecases

import (
	"errors"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type IReturnsUsecase interface {
	CreateNewReturn(userID, packageID uint) (*entities.Return, error)
	CreatePackageDispatch(userID, packageID uint) (*entities.Return, error)
	CreateReturnRequest(userID, packageID uint) (*entities.Return, error)
	CreatePackagePickup(userID, packageID uint) (*entities.PackagePickup, error)
}

type ReturnsUsecase struct {
	ReturnsRepo repositories.IReturnRepository
}

func (u ReturnsUsecase) CreateNewReturn(user *entities.User, packageID uint) (*entities.Return, error) {
	// var ret *entities.Return
	// var err error
	// if user.UserRole.SystemID == "restaurant" {
	// 	ret, err = u.CreatePackageDispatch(user.ID, packageID)
	// } else if user.UserRole.SystemID == "customer" {
	// 	ret, err = u.CreateReturnRequest(user.ID, packageID)
	// } else if user.UserRole.SystemID == "rider" {
	// 	ret, err = u.CreatePackagePickup(user.ID, packageID)
	// }

	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func (u ReturnsUsecase) CreatePackageDispatch(userID, packageID uint) (*entities.Return, error) {
	rets, err := u.ReturnsRepo.GetActiveReturns(packageID)
	if err != nil {
		return nil, err
	}

	if len(rets) > 0 {
		return nil, errors.New("")
	}

	ret, err := u.ReturnsRepo.CreateReturn(packageID)
	if err != nil {
		return nil, err
	}

	ret, err = u.ReturnsRepo.CreatePackageDispatch(ret.ID, userID)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (u ReturnsUsecase) CreateReturnRequest(userID, packageID, slotID uint) (*entities.Return, error) {
	rets, err := u.ReturnsRepo.GetActiveReturns(packageID)
	if err != nil {
		return nil, err
	}

	if len(rets) > 1 {
		return nil, errors.New("Invalid state; More then one open return found.")
	}

	if len(rets) == 1 && rets[0].Status != "Dispatched" {
		return nil, errors.New("No 'Dispatched' return found.")
	}

	var ret *entities.Return
	if len(rets) == 1 && rets[0].Status == "Dispatched" {
		ret = rets[0]
	} else if len(rets) == 0 {
		ret, err = u.ReturnsRepo.CreateReturn(packageID)
		if err != nil {
			return nil, err
		}
	}

	ret, err = u.ReturnsRepo.CreateReturnRequest(ret.ID, userID, slotID)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (u ReturnsUsecase) CreatePackagePickup(userID, packageID uint) (*entities.Return, error) {
	rets, err := u.ReturnsRepo.GetActiveReturns(packageID)
	if err != nil {
		return nil, err
	}

	if len(rets) == 0 {
		return nil, errors.New("No active return found.")
	}

	if len(rets) > 1 {
		return nil, errors.New("Invalid state; More then one open return found.")
	}

	if len(rets) == 1 && rets[0].Status != "Scheduled" {
		return nil, errors.New("No 'Scheduled' return found.")
	}

	if len(rets) == 1 && rets[0].Status == "Scheduled" {
		ret, err := u.ReturnsRepo.CreatePackageDispatch(rets[0].ID, userID)
		if err != nil {
			return nil, err
		}

		return ret, nil
	}

	return nil, errors.New("Unexpected; No conditions met")
}
