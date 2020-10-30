package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/usecases"
)

type UserChipsHandler struct {
	DB *gorm.DB
}

func (h UserChipsHandler) GetUserChips(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequestContext(r, w)
	if user == nil {
		return
	}

	urr := repositories.UserReturnRepository{DB: h.DB}
	u := usecases.UserChipsUsecase{
		UserReturnRepo: urr,
	}

	chips, err := u.GetUserChips(user.ID)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: "Unable to get the users chips",
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
	}

	writeSuccesResponse(chips, w)
}
