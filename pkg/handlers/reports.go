package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type ReportsHandler struct {
	DB *gorm.DB
}


func (h *ReportsHandler) StatusesPerRestaurant(w http.ResponseWriter, r *http.Request) {
	rr := repositories.ReportsRepository{h.DB}
	s, err := rr.InventoryPerStatusPerRestaurant()
	if err != nil {
		writeErrorResponse([]ResponseError{
			{
				Message: err.Error(),
				Code:    "",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(s, w)
}
