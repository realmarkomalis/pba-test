package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type PickupSlotsHandler struct {
	DB *gorm.DB
}

func (h PickupSlotsHandler) GetPickupSlots(w http.ResponseWriter, r *http.Request) {
	pr := repositories.PickupSlotsRepository{DB: h.DB}

	slots, err := pr.GetPickupSlots()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeSuccesResponse(slots, w)
}
