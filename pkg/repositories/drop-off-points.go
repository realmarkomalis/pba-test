package repositories

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type DropOffPointsRepository struct {
	DB *gorm.DB
}

func (r DropOffPointsRepository) GetDropOffPointByID(dopID uint) (*entities.DropOffPoint, error) {
	dop := models.DropOffPoint{}
	err := r.DB.
		Where("id = ?", dopID).
		First(&dop).
		Error
	if err != nil {
		return nil, err
	}

	d := dop.ModelToEntity()
	return &d, nil
}

func (r DropOffPointsRepository) ListAllDropOffPoints() ([]entities.DropOffPoint, error) {
	dos := []models.DropOffPoint{}
	err := r.DB.
		Preload("PackageTypes").
		Find(&dos).
		Error
	if err != nil {
		return nil, err
	}

	dops := []entities.DropOffPoint{}
	for _, d := range dos {
		dops = append(dops, d.ModelToEntity())
	}

	return dops, nil
}
