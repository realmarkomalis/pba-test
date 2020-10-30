package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type RestaurantsHandler struct {
	DB *gorm.DB
}

func (h RestaurantsHandler) GetRestaurants(w http.ResponseWriter, r *http.Request) {
	rr := repositories.RestaurantsRepository{DB: h.DB}

	restaurants, err := rr.GetRestaurants()
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(restaurants, w)
	return
}
