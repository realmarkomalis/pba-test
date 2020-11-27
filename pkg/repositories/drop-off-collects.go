package repositories

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type DropOffCollectsRepository struct {
	DB *gorm.DB
}

func (r DropOffCollectsRepository) CreateDropOffCollect(userID, returnID, dropOffPointID uint) (*entities.DropOffCollect, error) {
	user := models.User{}
	dop := models.DropOffPoint{}

	doi := models.DropOffCollect{
		UserID:         userID,
		ReturnID:       returnID,
		DropOffPointID: dropOffPointID,
	}
	err := r.DB.
		Create(&doi).
		Related(&user).
		Related(&dop).
		Error
	if err != nil {
		return nil, err
	}

	d := doi.ModelToEntity()
	return &d, nil
}
