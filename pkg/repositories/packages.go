package repositories

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type PackagesRepository struct {
	DB *gorm.DB
}

func (r PackagesRepository) GetPackage(packageID uint) (*entities.Package, error) {
	p := models.Package{}
	err := r.DB.
		Where("id = ?", packageID).
		Preload("PackageType").
		First(&p).
		Error
	if err != nil {
		return nil, err
	}

	d := p.ModelToEntity()
	return &d, nil
}
