package models

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type Restaurant struct {
	gorm.Model
	Name              string
	User              User
	UserID            uint
	DropOffPoint      DropOffPoint
	DropOffPointID    uint
	Latitude          float64
	Longitude         float64
	PostalCode        string
	StreetName        string
	HouseNumber       string
	HouseNumberSuffix string
	City              string
	Country           string
}

func (r Restaurant) ModelToEntity() entities.Restaurant {
	return entities.Restaurant{
		ID:                r.ID,
		Name:              r.Name,
		User:              r.User.ModelToEntity(),
		DropOffPoint:      r.DropOffPoint.ModelToEntity(),
		Latitude:          r.Latitude,
		Longitude:         r.Longitude,
		PostalCode:        r.PostalCode,
		StreetName:        r.StreetName,
		HouseNumber:       r.HouseNumber,
		HouseNumberSuffix: r.HouseNumberSuffix,
		City:              r.City,
		Country:           r.Country,
	}
}
