package repositories

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

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

	ps := models.PackageSupply{
		UserID:       userID,
		ReturnID:     ret.ID,
		RestaurantID: restaurantID,
	}
	err = r.DB.
		Create(&ps).
		Error
	if err != nil {
		return nil, err
	}

	err = r.DB.
		Model(&ret).
		Updates(models.Return{Status: models.Created}).
		Error
	if err != nil {
		return nil, err
	}

	ret.PackageSupply = ps
	return &entities.Return{
		ID:        ret.ID,
		Status:    ret.Status.String(),
		Package:   ret.Package.ModelToEntity(),
		CreatedAt: ret.CreatedAt,
	}, nil
}

func (r ReturnRepository) GetAllPackageSupplies() ([]*entities.Return, error) {
	rets := []models.Return{}
	supplies := []models.PackageSupply{}

	err := r.DB.
		Find(&supplies).
		Error
	if err != nil {
		return nil, err
	}

	returnIDs := []uint{}
	for _, d := range supplies {
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
		r := ret.ModelToEntity()
		returns = append(returns, &r)
	}

	return returns, nil
}
