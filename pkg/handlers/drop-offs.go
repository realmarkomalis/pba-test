package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type DropOffsHandler struct {
	DB *gorm.DB
}

func (h DropOffsHandler) CreateDropOffIntent(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	if u == nil {
		return
	}

	b := createDropOffIntentRequestBody{}
	if !validateRequestBody(r, w, &b) {
		return
	}

	dir := repositories.DropOffIntentsRepository{DB: h.DB}
	di, err := dir.CreateDropOffIntent(u.ID, b.PackageID, b.DropOffPointID)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusInternalServerError, w)
		return
	}

	writeSuccesResponse(di, w)
}
