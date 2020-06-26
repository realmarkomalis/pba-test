package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type Package struct {
	gorm.Model
	Name        string
	PackageCode uuid.UUID `gorm:"type:uuid;not null; default:uuid_generate_v4()"`
}

func (p Package) ModelToEntity() entities.Package {
	return entities.Package{
		ID:          p.ID,
		Name:        p.Name,
		PackageCode: p.PackageCode,
	}
}
