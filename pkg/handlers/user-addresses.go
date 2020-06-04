package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type UserAddressesHandler struct {
	DB *gorm.DB
}

type createAddressBody struct {
	FirstName         string `json:"first_name" valid:"type(string),stringlength(2|100)"`
	LastName          string `json:"last_name" valid:"type(string),stringlength(2|100)"`
	PostalCode        string `json:"postal_code" valid:"type(string),stringlength(6|10)"`
	StreetName        string `json:"street_name" valid:"type(string),stringlength(2|100)"`
	HouseNumber       string `json:"house_number" valid:"type(string),stringlength(1|10)"`
	HouseNumberSuffix string `json:"house_number_suffix" valid:"type(string),optional"`
	City              string `json:"city" valid:"type(string),stringlength(2|100)"`
}

func (h UserAddressesHandler) CreateUserAddress(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequestContext(r, w)
	pr := repositories.UserAddressesRepository{DB: h.DB}

	a := createAddressBody{}
	validateRequestBody(r, w, &a)

	slots, err := pr.CreateUserAddress(user, &entities.UserAddress{
		PostalCode:        a.PostalCode,
		StreetName:        a.StreetName,
		HouseNumber:       a.HouseNumber,
		HouseNumberSuffix: a.HouseNumberSuffix,
		City:              a.City,
		Country:           "NL",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeSuccesResponse(slots, w)
}

func (h UserAddressesHandler) UpdateUserAddress(w http.ResponseWriter, r *http.Request) {
	pr := repositories.PickupSlotsRepository{DB: h.DB}

	slots, err := pr.GetPickupSlots()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeSuccesResponse(slots, w)
}
