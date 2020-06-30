package repositories

import (
	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type IReturnRepository interface {
	GetReturn(returnID uint) (*entities.Return, error)
	GetReturns(packageID uint) ([]*entities.Return, error)
	GetActiveReturns(packageID uint) ([]*entities.Return, error)

	CreateReturn(packageID uint) (*entities.Return, error)
	CreatePackageSupply(returnID, userID, restaurantID uint) (*entities.Return, error)
	CreatePackageDispatch(returnID, userID uint) (*entities.Return, error)
	CreateReturnRequest(returnID, userID, slotID uint) (*entities.Return, error)
	CreatePackagePickup(returnID, userID uint) (*entities.Return, error)

	GetPackageDispatches(userID uint) ([]*entities.Return, error)
	GetPackagePickups(userID uint) ([]*entities.Return, error)
	GetReturnRequests(userID uint) ([]*entities.Return, error)
}

type ReturnRepository struct {
	DB *gorm.DB
}

func (r ReturnRepository) GetReturn(returnID uint) (*entities.Return, error) {
	rs := models.Return{}
	err := r.DB.
		Preload("Package").
		Where("id = ?", returnID).
		First(&rs).
		Error

	if err != nil {
		return nil, err
	}

	return &entities.Return{
		ID:     rs.ID,
		Status: rs.Status.String(),
		Package: entities.Package{
			ID:          rs.Package.ID,
			Name:        rs.Package.Name,
			PackageCode: rs.Package.PackageCode,
		},
		CreatedAt: rs.CreatedAt,
	}, nil
}

func (r ReturnRepository) GetReturns(packageID uint) ([]*entities.Return, error) {
	rs := []models.Return{}
	err := r.DB.
		Preload("Package").
		Where("package_id = ?", packageID).
		Find(&rs).
		Error

	if err != nil {
		return nil, err
	}

	rets := []*entities.Return{}
	for _, ret := range rs {
		rets = append(rets, &entities.Return{
			ID:     ret.ID,
			Status: ret.Status.String(),
			Package: entities.Package{
				ID:          ret.Package.ID,
				Name:        ret.Package.Name,
				PackageCode: ret.Package.PackageCode,
			},
			CreatedAt: ret.CreatedAt,
		})
	}

	return rets, nil
}

func (r ReturnRepository) GetActiveReturns(packageID uint) ([]*entities.Return, error) {
	rs := []models.Return{}
	err := r.DB.
		Where("package_id = ?", packageID).
		Not(&models.Return{Status: models.Fulfilled}).
		Find(&rs).
		Error

	if err != nil {
		return nil, err
	}

	rets := []*entities.Return{}
	for _, ret := range rs {
		rets = append(rets, &entities.Return{
			ID:     ret.ID,
			Status: ret.Status.String(),
			Package: entities.Package{
				ID:          ret.Package.ID,
				Name:        ret.Package.Name,
				PackageCode: ret.Package.PackageCode,
			},
			CreatedAt: ret.CreatedAt,
		})
	}

	return rets, nil
}

func (r ReturnRepository) CreateReturn(packageID uint) (*entities.Return, error) {
	pack := models.Package{}
	ret := models.Return{
		Status:    models.Created,
		PackageID: packageID,
	}
	err := r.DB.Create(&ret).Related(&pack).Error
	if err != nil {
		return nil, err
	}

	return &entities.Return{
		ID:     ret.ID,
		Status: ret.Status.String(),
		Package: entities.Package{
			ID:          pack.ID,
			Name:        pack.Name,
			PackageCode: pack.PackageCode,
		},
		CreatedAt: ret.CreatedAt,
	}, nil
}

func (r ReturnRepository) CreatePackageSupply(returnID, userID, restaurantID uint) (*entities.Return, error) {
	ret := models.Return{}
	err := r.DB.
		Preload("Package").
		Where("id = ?", returnID).
		First(&ret).
		Error

	if err != nil {
		return nil, err
	}

	err = r.DB.
		Model(&ret).
		Updates(models.Return{
			Status: models.Created,
			PackageSupply: models.PackageSupply{
				UserID:       userID,
				RestaurantID: restaurantID,
			},
		}).
		Error
	if err != nil {
		return nil, err
	}

	return &entities.Return{
		ID:        ret.ID,
		Status:    ret.Status.String(),
		Package:   ret.Package.ModelToEntity(),
		CreatedAt: ret.CreatedAt,
	}, nil
}

func (r ReturnRepository) CreatePackageDispatch(returnID, userID uint) (*entities.Return, error) {
	ret := models.Return{}
	err := r.DB.
		Preload("Package").
		Where("id = ?", returnID).
		First(&ret).
		Error

	if err != nil {
		return nil, err
	}

	err = r.DB.
		Model(&ret).
		Updates(models.Return{
			Status: models.Dispatched,
			PackageDispatch: models.PackageDispatch{
				UserID: userID,
			},
		}).
		Error
	if err != nil {
		return nil, err
	}

	return &entities.Return{
		ID:     ret.ID,
		Status: ret.Status.String(),
		Package: entities.Package{
			ID:          ret.Package.ID,
			Name:        ret.Package.Name,
			PackageCode: ret.Package.PackageCode,
		},
		CreatedAt: ret.CreatedAt,
	}, nil
}

func (r ReturnRepository) CreateReturnRequest(returnID, userID, slotID uint) (*entities.Return, error) {
	ret := models.Return{}
	err := r.DB.
		Preload("Package").
		Where("id = ?", returnID).
		First(&ret).
		Error

	if err != nil {
		return nil, err
	}

	err = r.DB.Model(&ret).Updates(models.Return{
		Status: models.Scheduled,
		ReturnRequest: models.ReturnRequest{
			UserID:       userID,
			ReturnID:     ret.ID,
			PickupSlotID: slotID,
		},
	}).Error
	if err != nil {
		return nil, err
	}

	return &entities.Return{
		ID:     ret.ID,
		Status: ret.Status.String(),
		Package: entities.Package{
			ID:          ret.Package.ID,
			Name:        ret.Package.Name,
			PackageCode: ret.Package.PackageCode,
		},
		CreatedAt: ret.CreatedAt,
	}, nil
}

func (r ReturnRepository) CreatePackagePickup(returnID, userID uint) (*entities.Return, error) {
	ret := models.Return{}
	err := r.DB.
		Preload("Package").
		Where("id = ?", returnID).
		First(&ret).
		Error

	if err != nil {
		return nil, err
	}

	err = r.DB.
		Model(&ret).
		Updates(models.Return{
			Status: models.Fulfilled,
			PackagePickup: models.PackagePickup{
				UserID: userID,
			},
		}).
		Error
	if err != nil {
		return nil, err
	}

	return &entities.Return{
		ID:     ret.ID,
		Status: ret.Status.String(),
		Package: entities.Package{
			ID:          ret.Package.ID,
			Name:        ret.Package.Name,
			PackageCode: ret.Package.PackageCode,
		},
		CreatedAt: ret.CreatedAt,
	}, nil
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
				ID:          ret.Package.ID,
				Name:        ret.Package.Name,
				PackageCode: ret.Package.PackageCode,
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
				ID:          ret.Package.ID,
				Name:        ret.Package.Name,
				PackageCode: ret.Package.PackageCode,
			},
		})
	}

	return returns, nil
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
				ID:          ret.Package.ID,
				Name:        ret.Package.Name,
				PackageCode: ret.Package.PackageCode,
			},
		})
	}

	return returns, nil
}

func (r ReturnRepository) GetAllReturnRequests() ([]entities.PickupRequest, error) {
	rets := []models.Return{}
	err := r.DB.
		Preload("ReturnRequest").
		Preload("Package").
		Where("status = ?", models.Scheduled).
		Find(&rets).
		Error
	if err != nil {
		return nil, err
	}

	prIDs := []uint{}
	for _, ret := range rets {
		prIDs = append(prIDs, ret.ReturnRequest.ID)
	}

	rrs := []models.ReturnRequest{}
	r.DB.Preload("PickupSlot").Preload("User").Preload("User.UserAddresses").Where(prIDs).Find(&rrs)

	prs := []entities.PickupRequest{}
	for _, ret := range rrs {
		address := entities.UserAddress{}
		if len(ret.User.UserAddresses) > 0 {
			a := ret.User.UserAddresses[0]
			address.PostalCode = a.PostalCode
			address.StreetName = a.StreetName
			address.HouseNumber = a.HouseNumber
			address.HouseNumberSuffix = a.HouseNumberSuffix
			address.City = a.City
		}
		prs = append(prs, entities.PickupRequest{
			ID: ret.ID,
			User: entities.User{
				ID:    ret.User.ID,
				Email: ret.User.Email,
				UserAddresses: []entities.UserAddress{
					address,
				},
			},
		})
		prIDs = append(prIDs, ret.ID)
	}

	return prs, nil
}
