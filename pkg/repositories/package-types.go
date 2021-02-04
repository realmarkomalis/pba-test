package repositories

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type PackageTypesRepository struct {
	DB *gorm.DB
}

func (r PackageTypesRepository) GetPackageTypes() ([]entities.PackageType, error) {
	pt := []models.PackageType{}
	err := r.DB.
		Find(&pt).
		Error
	if err != nil {
		return nil, err
	}

	ts := []entities.PackageType{}
	for _, t := range pt {
		ts = append(ts, t.ModelToEntity())
	}
	return ts, nil
}
