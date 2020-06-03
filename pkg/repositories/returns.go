package repositories

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type IReturnRepository interface {
	Create(userID, packageID uint) (*entities.Return, error)
	CreatePackageDispatch(userID, packageID uint) (*entities.Return, error)
	CreateReturnRequest(userID, packageID, slotID uint) (*entities.Return, error)
	CreatePackagePickup(userID, packageID uint) (*entities.PackagePickup, error)
	GetPackageDispatches(userID uint) ([]*entities.Return, error)
	GetPackagePickups(userID uint) ([]*entities.Return, error)
	GetReturnRequests(userID uint) ([]*entities.Return, error)
}

type ReturnRepository struct {
	DB *gorm.DB
}

func (r ReturnRepository) Create(userID, packageID uint) (*entities.Return, error) {
	ret := models.Return{}
	results := 0
	err := r.DB.
		Where("package_id = ?", packageID).
		Not(&models.Return{Status: models.Fulfilled}).
		Find(&ret).
		Count(&results).
		Error

	if results > 0 {
		fmt.Printf("found %d number of results", results)
		fmt.Println(err)
		return nil, fmt.Errorf("Already active return for package with ID: %d", packageID)
	}

	pack := models.Package{}
	ret = models.Return{
		Status:    models.Dispatched,
		PackageID: packageID,
		PackageDispatch: models.PackageDispatch{
			UserID: userID,
		},
	}
	err = r.DB.Create(&ret).Related(&pack).Error
	if err != nil {
		return nil, err
	}

	return &entities.Return{
		ID:     ret.ID,
		Status: models.Dispatched.String(),
		Package: entities.Package{
			ID:   pack.ID,
			Name: pack.Name,
		},
	}, nil
}

func (r ReturnRepository) CreatePackageDispatch(userID, packageID uint) (*entities.Return, error) {
	return nil, nil
}

func (r ReturnRepository) CreateReturnRequest(userID, packageID, slotID uint) (*entities.Return, error) {
	ret := models.Return{}
	err := r.DB.
		Preload("Package").
		Where("package_id = ? AND status = ?", packageID, models.Dispatched).
		First(&ret).
		Error
	if err != nil {
		fmt.Printf("error getting package %d: %s", packageID, err)
		return nil, err
	}
	if ret.Status != models.Dispatched {
		return nil, errors.New("Couldn't find an dispatched return for this package")
	}

	slot := models.PickupSlot{}
	err = r.DB.
		Where("id = ?", slotID).
		First(&slot).
		Error
	if err != nil {
		return nil, err
	}

	err = r.DB.Model(&slot).Updates(models.PickupSlot{
		UserID: userID,
		Booked: true,
	}).Error
	if err != nil {
		return nil, err
	}

	err = r.DB.Model(&ret).Updates(models.Return{
		Status: models.Scheduled,
		ReturnRequest: models.ReturnRequest{
			UserID:       userID,
			ReturnID:     ret.ID,
			PickupSlotID: slot.ID,
		},
	}).Error
	if err != nil {
		return nil, err
	}

	return &entities.Return{
		ID:     ret.ID,
		Status: ret.Status.String(),
		Package: entities.Package{
			ID:   ret.Package.ID,
			Name: ret.Package.Name,
		},
	}, nil
}

func (r ReturnRepository) CreatePackagePickup(userID, packageID uint) (*entities.PackagePickup, error) {
	ret := models.Return{}
	err := r.DB.
		Preload("Package").
		Where("package_id = ? AND status = ?", packageID, models.Scheduled).
		First(&ret).
		Error
	if err != nil {
		return nil, err
	}
	if ret.Status != models.Scheduled {
		return nil, errors.New("Couldn't find an scheduled return for this package")
	}

	err = r.DB.Model(&ret).Updates(models.Return{
		Status: models.Fulfilled,
		PackagePickup: models.PackagePickup{
			UserID: userID,
		},
	}).Error
	if err != nil {
		return nil, err
	}

	return &entities.PackagePickup{
		ID: ret.PackagePickup.ID,
		Return: entities.Return{
			ID:     ret.ID,
			Status: ret.Status.String(),
			Package: entities.Package{
				ID:   ret.Package.ID,
				Name: ret.Package.Name,
			},
		},
	}, nil
}

func (r ReturnRepository) GetPackagePickups(userID uint) ([]*entities.Return, error) {
	rets := []models.Return{}
	pickups := []models.PackagePickup{}

	err := r.DB.Where("user_id = ?", userID).Find(&pickups).Error
	if err != nil {
		return nil, err
	}

	returnIDs := []uint{}
	for _, p := range pickups {
		returnIDs = append(returnIDs, p.ReturnID)
	}

	err = r.DB.
		Preload("Package").
		Where("id in (?)", returnIDs).
		Find(&rets).
		Error
	if err != nil {
		return nil, err
	}

	returns := []*entities.Return{}
	for _, ret := range rets {
		returns = append(returns, &entities.Return{
			ID:        ret.ID,
			Status:    ret.Status.String(),
			CreatedAt: ret.CreatedAt,
			Package: entities.Package{
				ID:   ret.Package.ID,
				Name: ret.Package.Name,
			},
		})
	}

	return returns, nil
}

func (r ReturnRepository) GetPackageDispatches(userID uint) ([]*entities.Return, error) {
	rets := []models.Return{}
	dispatches := []models.PackageDispatch{}

	err := r.DB.
		Where("user_id = ?", userID).
		Find(&dispatches).
		Error
	if err != nil {
		return nil, err
	}

	returnIDs := []uint{}
	for _, d := range dispatches {
		returnIDs = append(returnIDs, d.ReturnID)
	}

	err = r.DB.
		Preload("Package").
		Where("id in (?)", returnIDs).
		Find(&rets).
		Error
	if err != nil {
		return nil, err
	}

	returns := []*entities.Return{}
	for _, ret := range rets {
		returns = append(returns, &entities.Return{
			ID:        ret.ID,
			Status:    ret.Status.String(),
			CreatedAt: ret.CreatedAt,
			Package: entities.Package{
				ID:   ret.Package.ID,
				Name: ret.Package.Name,
			},
		})
	}

	return returns, nil
}

func (r ReturnRepository) GetReturnRequests(userID uint) ([]*entities.Return, error) {
	rets := []models.Return{}
	returnRequests := []models.ReturnRequest{}

	err := r.DB.
		Where("user_id = ?", userID).
		Find(&returnRequests).
		Error
	if err != nil {
		return nil, err
	}

	returnIDs := []uint{}
	for _, d := range returnRequests {
		returnIDs = append(returnIDs, d.ReturnID)
	}

	err = r.DB.
		Preload("Package").
		Where("id in (?)", returnIDs).
		Find(&rets).
		Error
	if err != nil {
		return nil, err
	}

	returns := []*entities.Return{}
	for _, ret := range rets {
		returns = append(returns, &entities.Return{
			ID:        ret.ID,
			Status:    ret.Status.String(),
			CreatedAt: ret.CreatedAt,
			Package: entities.Package{
				ID:   ret.Package.ID,
				Name: ret.Package.Name,
			},
		})
	}

	return returns, nil
}
