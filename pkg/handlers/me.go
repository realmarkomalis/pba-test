package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/services"
)

type MeHandler struct {
	DB *gorm.DB
}

type myQRCodeResponse struct {
	Code string `json:"code"`
}

type updateUserRequestBody struct {
	FirstName   string `json:"first_name" valid:"type(string),optional"`
	LastName    string `json:"last_name" valid:"type(string),optional"`
	PhoneNumber string `json:"phonenumber" valid:"type(string),optional"`
	IBAN        string `json:"iban" valid:"type(string),optional"`
}

func (h MeHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	ur := repositories.UserRepository{DB: h.DB}

	user, err := ur.GetByEmail(strings.ToLower(u.Email))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ts := services.NewTokenService(
		viper.GetString("token_service.secret_key"),
		viper.GetString("token_service.issuer"),
		viper.GetInt("token_service.minutes_valid"),
	)
	token, err := ts.Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Token = token

	writeSuccesResponse(user, w)
}

func (h MeHandler) GetMyQR(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)

	s := services.QRCodeService{}
	code, err := s.CreateBase64Code(fmt.Sprintf("%d", u.ID))
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusInternalServerError, w)
		return
	}

	writeSuccesResponse(myQRCodeResponse{Code: code}, w)
}

func (h MeHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)

	b := updateUserRequestBody{}
	if !validateRequestBody(r, w, &b) {
		return
	}

	ur := repositories.UserRepository{DB: h.DB}
	user, err := ur.UpdateUser(&entities.User{
		ID:          u.ID,
		FirstName:   b.FirstName,
		LastName:    b.LastName,
		PhoneNumber: b.PhoneNumber,
		IBAN:        b.IBAN,
	})

	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(user, w)
}
