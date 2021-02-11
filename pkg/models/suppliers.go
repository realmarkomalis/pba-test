package models

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type Supplier struct {
	gorm.Model
	Name         string
	Users        []User `gorm:"many2many:suppplier_users;"`
	PackageTypes []PackageType
}

func (s Supplier) ModelToEntity() entities.Supplier {
	return entities.Supplier{
		ID:   s.ID,
		Name: s.Name,
	}
}
