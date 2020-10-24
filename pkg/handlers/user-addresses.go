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
	PhoneNumber       string `json:"phonenumber" valid:"type(string),stringlength(8|20),optional"`
}

func (h UserAddressesHandler) CreateUserAddress(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequestContext(r, w)

	a := createAddressBody{}
	if !validateRequestBody(r, w, &a) {
		return
	}

	pr := repositories.UserAddressesRepository{DB: h.DB}
	addr, err := pr.UpdateUserAddress(user, &entities.UserAddress{
		PostalCode:        a.PostalCode,
		StreetName:        a.StreetName,
		HouseNumber:       a.HouseNumber,
		HouseNumberSuffix: a.HouseNumberSuffix,
		City:              a.City,
		Country:           "NL",
	})
	if err != nil {
		writeErrorResponse([]ResponseError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	ur := repositories.UserRepository{DB: h.DB}
	_, err = ur.UpdateUser(&entities.User{
		ID:          user.ID,
		FirstName:   a.FirstName,
		LastName:    a.LastName,
		PhoneNumber: a.PhoneNumber,
	})
	if err != nil {
		writeErrorResponse([]ResponseError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(addr, w)
}

func (h UserAddressesHandler) UpdateUserAddress(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequestContext(r, w)

	a := createAddressBody{}
	if !validateRequestBody(r, w, &a) {
		return
	}

	pr := repositories.UserAddressesRepository{DB: h.DB}
	addr, err := pr.UpdateUserAddress(user, &entities.UserAddress{
		PostalCode:        a.PostalCode,
		StreetName:        a.StreetName,
		HouseNumber:       a.HouseNumber,
		HouseNumberSuffix: a.HouseNumberSuffix,
		City:              a.City,
		Country:           "NL",
	})
	if err != nil {
		writeErrorResponse([]ResponseError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	ur := repositories.UserRepository{DB: h.DB}
	_, err = ur.UpdateUser(&entities.User{
		ID:        user.ID,
		FirstName: a.FirstName,
		LastName:  a.LastName,
	})
	if err != nil {
		writeErrorResponse([]ResponseError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(addr, w)
}
