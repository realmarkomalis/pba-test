package repositories

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type PackagesRepository struct {
	DB *gorm.DB
}

func (r PackagesRepository) CreateOrUpdatePackage(code string, packageType uint) (*entities.Package, error) {
	p := models.Package{}

	u, err := uuid.FromString(code)
	if err != nil {
		return nil, err
	}

	err = r.DB.
		Where(models.Package{PackageCode: u}).
		Assign(models.Package{PackageTypeID: packageType}).
		FirstOrCreate(&p).
		Error
	if err != nil {
		return nil, err
	}

	pa := p.ModelToEntity()
	return &pa, nil
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

func (r PackagesRepository) GetPackageByCode(code string) (*entities.Package, error) {
	p := models.Package{}
	err := r.DB.
		Where("package_code = ?", code).
		Preload("PackageType").
		First(&p).
		Error
	if err != nil {
		return nil, err
	}

	d := p.ModelToEntity()
	return &d, nil
}
