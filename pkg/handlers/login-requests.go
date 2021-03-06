package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/services"
	"gitlab.com/markomalis/packback-api/pkg/usecases"
)

type LoginRequestsHandler struct {
	DB *gorm.DB
}

type createRequestBody struct {
	Email string `json:"email" valid:"email"`
}

type requestLoginResponse struct {
	Requested bool `json:"requested"`
}

func (h LoginRequestsHandler) Create(w http.ResponseWriter, r *http.Request) {
	ur := repositories.UserRepository{h.DB}
	as := services.PasswordlessAuth{
		Domain:       viper.GetString("auth_service.domain"),
		ClientID:     viper.GetString("auth_service.client_id"),
		ClientSecret: viper.GetString("auth_service.client_secret"),
		RedirectURL:  viper.GetString("auth_service.redirect_url"),
	}
	um := services.UserManagementService{
		Domain:       viper.GetString("user_management_service.domain"),
		ClientID:     viper.GetString("user_management_service.client_id"),
		ClientSecret: viper.GetString("user_management_service.client_secret"),
		Audience:     viper.GetString("user_management_service.audience"),
		AccessToken:  "",
	}
	ts := services.NewTokenService(
		viper.GetString("token_service.secret_key"),
		viper.GetString("token_service.issuer"),
		viper.GetInt("token_service.minutes_valid"),
	)
	a := usecases.AuthUsecase{
		ur,
		as,
		um,
		*ts,
	}

	body := createRequestBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := govalidator.ValidateStruct(body)
	if err != nil || !result {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requested, err := a.RequestLogin(strings.ToLower(body.Email))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(&requestLoginResponse{requested})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
