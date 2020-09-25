package entities

type User struct {
	ID            uint          `json:"id"`
	Email         string        `json:"email"`
	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	PhoneNumber   string        `json:"phonenumber"`
	UserAddresses []UserAddress `json:"user_addresses"`
	UserRole      UserRole      `json:"user_role"`
	Token         string        `json:"token"`
}

type UserAddress struct {
	ID                uint   `json:"id"`
	PostalCode        string `json:"postal_code"`
	StreetName        string `json:"street_name"`
	HouseNumber       string `json:"house_number"`
	HouseNumberSuffix string `json:"house_number_suffix"`
	City              string `json:"city"`
	Country           string `json:"country"`
}

type UserRole struct {
	Name     string `json:"name"`
	SystemID string `json:"system_id"`
}
