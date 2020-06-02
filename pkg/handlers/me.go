package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type MeHandler struct {
	DB *gorm.DB
}

func (h MeHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	ur := repositories.UserRepository{DB: h.DB}

	user, err := ur.GetByEmail(u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeSuccesResponse(user, w)
}
