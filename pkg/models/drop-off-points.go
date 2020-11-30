package models

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type DropOffPoint struct {
	gorm.Model
	Name              string  `gorm:"not null"`
	Latitude          float64 `gorm:"not null"`
	Longitude         float64 `gorm:"not null"`
	PostalCode        string  `gorm:"not null"`
	StreetName        string  `gorm:"not null"`
	HouseNumber       string  `gorm:"not null"`
	HouseNumberSuffix string
	City              string        `gorm:"not null"`
	Country           string        `gorm:"not null"`
	PackageTypes      []PackageType `gorm:"many2many:drop_off_point_package_types;"`
}

func (d DropOffPoint) ModelToEntity() entities.DropOffPoint {
	pt := []entities.PackageType{}
	for _, p := range d.PackageTypes {
		pt = append(pt, p.ModelToEntity())
	}

	return entities.DropOffPoint{
		ID:                d.ID,
		Name:              d.Name,
		Latitude:          d.Latitude,
		Longitude:         d.Longitude,
		PostalCode:        d.PostalCode,
		StreetName:        d.StreetName,
		HouseNumber:       d.HouseNumber,
		HouseNumberSuffix: d.HouseNumberSuffix,
		City:              d.City,
		Country:           d.Country,
		PackageTypes:      pt,
	}
}
