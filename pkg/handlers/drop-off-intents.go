package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/usecases"
)

type DropOffIntentsHandler struct {
	DB *gorm.DB
}

type createDropOffIntentRequestBody struct {
	PackageID      uint `json:"package_id" valid:"required"`
	DropOffPointID uint `json:"drop_off_point_id" valid:"optional"`
}

func (h DropOffIntentsHandler) CreateDropOffIntent(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	if u == nil {
		return
	}

	b := createDropOffIntentRequestBody{}
	if !validateRequestBody(r, w, &b) {
		return
	}

	rr := repositories.ReturnRepository{DB: h.DB}
	dr := repositories.DropOffIntentsRepository{DB: h.DB}
	du := usecases.DropOffIntentUsecases{
		ReturnsRepo:        rr,
		DropOffIntentsRepo: dr,
	}

	di, err := du.CreateDropOffIntent(u.ID, b.PackageID, b.DropOffPointID)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(di, w)
}

func (h DropOffIntentsHandler) ListDropOffIntents(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	if u == nil {
		return
	}

	dr := repositories.DropOffIntentsRepository{DB: h.DB}
	di, err := dr.GetUserDropOffIntents(u.ID)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(di, w)
}
