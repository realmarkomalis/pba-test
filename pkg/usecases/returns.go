package usecases

import (
	"errors"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type IReturnsUsecase interface {
	CreatePackageSupplies(userID uint, packageIDs []uint, restaurantID uint) ([]*entities.Return, error)
	CreatePackageSupply(userID uint, packageID uint, restaurantID uint) (*entities.Return, error)

	CreatePackageDispatch(userID, packageID uint) (*entities.Return, error)
	CreatePackageDispatches(userID uint, packageIDs []uint) ([]*entities.Return, error)

	CreateReturnRequest(userID, packageID, slotID uint, comment string) (*entities.Return, error)
	CreateReturnRequests(userID, slotID uint, packageIDs []uint, comment string) ([]*entities.Return, error)

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

func (u ReturnsUsecase) CreatePackageSupplies(userID uint, packageIDs []uint, restaurantID uint) ([]*entities.Return, error) {
	rets := []*entities.Return{}
	returnIDs := []uint{}
	for _, id := range packageIDs {
		ret, err := u.CreatePackageSupply(userID, id, restaurantID)
		if err == nil {
			rets = append(rets, ret)
			returnIDs = append(returnIDs, ret.ID)
		}
	}

	if len(rets) == 0 {
		return rets, nil
	}

	_, err := u.UserReturnRepo.CreateUserReturnEntry(userID, returnIDs)
	if err != nil {
		return nil, err
	}

	return rets, nil
}

func (u ReturnsUsecase) CreatePackageSupply(userID uint, packageID uint, restaurantID uint) (*entities.Return, error) {
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

	ret, err = u.ReturnsRepo.CreatePackageSupply(ret.ID, userID, restaurantID)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

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

	if len(rets) == 0 {
		return rets, nil
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

	if len(rets) == 0 {
		return nil, errors.New("No supplied package found")
	}

	if len(rets) > 1 {
		return nil, errors.New("Inconsitent state: Multiple active returns found")
	}

	if len(rets) == 1 && rets[0].Status != "Created" {
		return nil, errors.New("No supplied packages found")
	}

	ret := rets[0]
	ret, err = u.ReturnsRepo.CreatePackageDispatch(ret.ID, userID)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (u ReturnsUsecase) CreateReturnRequests(userID, slotID uint, packageIDs []uint, comment string) ([]*entities.Return, error) {
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
		ret, err := u.CreateReturnRequest(userID, id, slot.ID, comment)
		if err == nil {
			rets = append(rets, ret)
			returnIDs = append(returnIDs, ret.ID)
		}
	}

	if len(rets) == 0 {
		return rets, nil
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

func (u ReturnsUsecase) CreateReturnRequest(userID, packageID, slotID uint, comment string) (*entities.Return, error) {
	rets, err := u.ReturnsRepo.GetActiveReturns(packageID)
	if err != nil {
		return nil, err
	}

	if len(rets) == 0 {
		return nil, errors.New("No active package found")
	}

	if len(rets) > 1 {
		return nil, errors.New("Invalid state; More then one open return found.")
	}

	if len(rets) == 1 && (rets[0].Status != "Dispatched" && rets[0].Status != "Created") {
		return nil, errors.New("No 'Created' or 'Dispatched'` return found.")
	}

	ret := rets[0]
	ret, err = u.ReturnsRepo.CreateReturnRequest(ret.ID, userID, slotID, comment)
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

	if len(rets) == 0 {
		return rets, nil
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
		ret, err := u.ReturnsRepo.CreatePackagePickup(rets[0].ID, userID)
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
