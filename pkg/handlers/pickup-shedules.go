package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/usecases"
)

type PickupShedulesHandler struct {
	DB *gorm.DB
}

func (h PickupShedulesHandler) GetPickupShedules(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequestContext(r, w)
	if user.UserRole.SystemID != "admin" {
		http.Error(w, "Invalid user type for this request", http.StatusNotAcceptable)
		return
	}

	rr := repositories.UserReturnRepository{DB: h.DB}
	u := usecases.PickupScheduleUsecase{rr}

	rets, err := u.GetScheduledReturnEntries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeSuccesResponse(rets, w)
}
