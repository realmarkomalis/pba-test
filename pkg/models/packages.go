package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type Package struct {
	gorm.Model
	Name          string
	PackageCode   uuid.UUID `gorm:"type:uuid;not null; default:uuid_generate_v4()"`
	PackageTypeID uint
	PackageType   PackageType
}

func (p Package) ModelToEntity() entities.Package {
	return entities.Package{
		ID:          p.ID,
		Name:        p.Name,
		PackageCode: p.PackageCode,
		PackageType: p.PackageType.ModelToEntity(),
	}
}

type PackageType struct {
	gorm.Model
	Name          string
	Description   string
	Manufacturer  string
	ProductCode   string
	DepositAmount float64
}

func (p PackageType) ModelToEntity() entities.PackageType {
	return entities.PackageType{
		ID:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		Manufacturer:  p.Manufacturer,
		ProductCode:   p.ProductCode,
		DepositAmount: p.DepositAmount,
	}
}
