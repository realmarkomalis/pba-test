package usecases

import (
	"errors"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type IReturnsUsecase interface {
	CreatePackageDispatch(userID, packageID uint) (*entities.Return, error)
	CreatePackageDispatches(userID uint, packageIDs []uint) ([]*entities.Return, error)

	CreateReturnRequest(userID, packageID, slotID uint) (*entities.Return, error)
	CreateReturnRequests(userID, slotID uint, packageIDs []uint) ([]*entities.Return, error)

	CreatePackagePickup(userID, packageID uint) (*entities.PackagePickup, error)
	CreatePackagePickups(userID uint, packageIDs []uint) ([]*entities.Return, error)

	GetUserReturns(userID uint) ([]*entities.UserReturnEntry, error)
	GetUserReturnEntries(userID uint) ([]*entities.UserReturnEntry, error)
}

type ReturnsUsecase struct {
	ReturnsRepo     repositories.IReturnRepository
	UserReturnRepo  repositories.IUserReturnRepository
	PickupSlotsRepo repositories.PickupSlotsRepository
}

type Bla map[string]interface{}

func (u ReturnsUsecase) CreatePackageDispatches(userID uint, packageIDs []uint) ([]*entities.Return, error) {
	rets := []*entities.Return{}
	returnIDs := []uint{}
	for _, id := range packageIDs {
		ret, err := u.CreatePackageDispatch(userID, id)
		if err == nil {
			rets = append(rets, ret)
			returnIDs = append(returnIDs, ret.ID)
		}
	}

	_, err := u.UserReturnRepo.CreateUserReturnEntry(userID, returnIDs)
	if err != nil {
		return nil, err
	}

	return rets, nil
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

func (u ReturnsUsecase) CreateReturnRequests(userID, slotID uint, packageIDs []uint) ([]*entities.Return, error) {
	slot, err := u.PickupSlotsRepo.GetPickupSlot(slotID)
	if err != nil {
		return nil, err
	}

	if slot.Booked {
		return nil, errors.New("Provided slot is already booked")
	}

	rets := []*entities.Return{}
	returnIDs := []uint{}
	for _, id := range packageIDs {
		ret, err := u.CreateReturnRequest(userID, id, slot.ID)
		if err == nil {
			rets = append(rets, ret)
			returnIDs = append(returnIDs, ret.ID)
		}
	}

	_, err = u.PickupSlotsRepo.BookPickupSlot(slot.ID)
	if err != nil {
		return nil, err
	}

	_, err = u.UserReturnRepo.CreateUserReturnEntry(userID, returnIDs)
	if err != nil {
		return nil, err
	}

	return rets, nil
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

func (u ReturnsUsecase) CreatePackagePickups(userID uint, packageIDs []uint) ([]*entities.Return, error) {
	rets := []*entities.Return{}
	returnIDs := []uint{}
	for _, id := range packageIDs {
		ret, err := u.CreatePackagePickup(userID, id)
		if err == nil {
			rets = append(rets, ret)
			returnIDs = append(returnIDs, ret.ID)
		}
	}

	_, err := u.UserReturnRepo.CreateUserReturnEntry(userID, returnIDs)
	if err != nil {
		return nil, err
	}

	return rets, nil
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

func (u ReturnsUsecase) GetUserReturns(userID uint) ([]*entities.UserReturnEntry, error) {
	return u.GetUserReturnEntries(userID)
}

func (u ReturnsUsecase) GetUserReturnEntries(userID uint) ([]*entities.UserReturnEntry, error) {
	return u.UserReturnRepo.GetUserReturnEntries(userID)
}
