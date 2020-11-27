package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/usecases"
)

type DropOffCollectsHandler struct {
	DB *gorm.DB
}

type createDropOffCollectRequestBody struct {
	PackageID      uint `json:"package_id" valid:"required"`
	DropOffPointID uint `json:"drop_off_point_id" valid:"required"`
}

func (h DropOffCollectsHandler) CreateDropOffCollect(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	if u == nil {
		writeErrorResponse([]entities.APIError{}, http.StatusForbidden, w)
		return
	}

	b := createDropOffCollectRequestBody{}
	if !validateRequestBody(r, w, &b) {
		writeErrorResponse([]entities.APIError{}, http.StatusInternalServerError, w)
		return
	}

	du := usecases.DropOffCollectUsecases{
		ReturnsRepository:         repositories.ReturnRepository{DB: h.DB},
		DropOffCollectsRepository: repositories.DropOffCollectsRepository{DB: h.DB},
		RestaurantsRepository:     repositories.RestaurantsRepository{DB: h.DB},
	}

	var dc *entities.DropOffCollect
	var err error
	if u.UserRole.Name == "restaurant" {
		dc, err = du.CreateDropOffCollectAndSupply(u.ID, b.PackageID, b.DropOffPointID)
	} else if u.UserRole.Name == "package_supplier" {
		dc, err = du.CreateDropOffCollect(u.ID, b.PackageID, b.DropOffPointID)
	}

	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(dc, w)
}
