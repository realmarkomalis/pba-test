package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type DropOffPointsHandler struct {
	DB *gorm.DB
}

func (h DropOffPointsHandler) GetDropOffPoint(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["dop_id"], 10, 64)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}
	dr := repositories.DropOffPointsRepository{DB: h.DB}
	dop, err := dr.GetDropOffPointByID(uint(id))
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusInternalServerError, w)
		return
	}

	writeSuccesResponse(dop, w)
}

func (h DropOffPointsHandler) ListDropOffPoints(w http.ResponseWriter, r *http.Request) {
	dr := repositories.DropOffPointsRepository{DB: h.DB}
	dops, err := dr.ListAllDropOffPoints()
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusInternalServerError, w)
		return
	}

	writeSuccesResponse(dops, w)
}
