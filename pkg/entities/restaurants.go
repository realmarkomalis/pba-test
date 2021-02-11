package entities

type Restaurant struct {
	ID                uint         `json:"id"`
	Name              string       `json:"name"`
	User              User         `json:"user"`
	DropOffPoint      DropOffPoint `json:"drop_off_point"`
	Latitude          float64      `json:"latitude"`
	Longitude         float64      `json:"longitude"`
	PostalCode        string       `json:"postal_code"`
	StreetName        string       `json:"street_name"`
	HouseNumber       string       `json:"house_number"`
	HouseNumberSuffix string       `json:"house_number_suffix"`
	City              string       `json:"city"`
	Country           string       `json:"country"`
	Deactivated       bool         `json:"deactivated"`
}
