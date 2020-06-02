package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/services"
	"gitlab.com/markomalis/packback-api/pkg/usecases"
)

type LoginAttemptsHandler struct {
	DB *gorm.DB
}

type loginAttemptBody struct {
	Code  string `json:"code" valid:"type(string)"`
	Email string `json:"email" valid:"email"`
}

func (h *LoginAttemptsHandler) Create(w http.ResponseWriter, r *http.Request) {
	ur := repositories.UserRepository{h.DB}
	as := services.PasswordlessAuth{
		Domain:       viper.GetString("auth_service.domain"),
		ClientID:     viper.GetString("auth_service.client_id"),
		ClientSecret: viper.GetString("auth_service.client_secret"),
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

	body := loginAttemptBody{}
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

	user, err := a.LoginUserWithCode(body.Email, body.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
