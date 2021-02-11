package entities

import "time"

type DropOffPoint struct {
	ID                uint          `json:"id"`
	Name              string        `json:"name"`
	Latitude          float64       `json:"latitude"`
	Longitude         float64       `json:"longitude"`
	PostalCode        string        `json:"postal_code"`
	StreetName        string        `json:"street_name"`
	HouseNumber       string        `json:"house_number"`
	HouseNumberSuffix string        `json:"house_number_suffix"`
	City              string        `json:"city"`
	Country           string        `json:"country"`
	Deactivated       bool          `json:"deactivated"`
	PackageTypes      []PackageType `json:"package_types"`
}

type DropOffIntent struct {
	ID           uint         `json:"id"`
	CreatedAt    time.Time    `json:"created_at"`
	User         User         `json:"user"`
	ReturnID     uint         `json:"return_id"`
	DropOffPoint DropOffPoint `json:"drop_off_point"`
}

type DropOffCollect struct {
	ID           uint         `json:"id"`
	CreatedAt    time.Time    `json:"created_at"`
	ReturnID     uint         `json:"return_id"`
	User         User         `json:"user"`
	DropOffPoint DropOffPoint `json:"drop_off_point"`
}
