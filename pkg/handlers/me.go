package handlers

import (
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/services"
)

type MeHandler struct {
	DB *gorm.DB
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
