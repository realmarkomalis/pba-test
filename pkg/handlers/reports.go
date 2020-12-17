package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/usecases"
)

type ReportsHandler struct {
	DB *gorm.DB
}

func (h *ReportsHandler) ListReportableRestaurants(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequestContext(r, w)
	if user.UserRole.Name != "admin" && user.UserRole.Name != "package_supplier" && user.UserRole.Name != "restaurant" && user.UserRole.Name != "closed_loop_restaurant" {
		writeErrorResponse([]entities.APIError{}, http.StatusForbidden, w)
		return
	}

	u := usecases.ReportsUsecase{
		ReportsRepository:     repositories.ReportsRepository{DB: h.DB},
		RestaurantsRepository: repositories.RestaurantsRepository{DB: h.DB},
	}

	if user.UserRole.Name == "admin" {
		rs, err := u.GetAdminRestaurants()
		if err != nil {
			writeErrorResponse([]entities.APIError{
				{
					Message: err.Error(),
					Code:    "",
				},
			}, http.StatusBadRequest, w)
			return
		}

		writeSuccesResponse(rs, w)
		return
	}

	if user.UserRole.Name == "package_supplier" {
		rs, err := u.GetRepotableRestaurants(user.ID)
		if err != nil {
			writeErrorResponse([]entities.APIError{
				{
					Message: err.Error(),
					Code:    "",
				},
			}, http.StatusBadRequest, w)
			return
		}

		writeSuccesResponse(rs, w)
		return
	}

	if user.UserRole.Name == "restaurant" || user.UserRole.Name == "closed_loop_restaurant" {
		rs, err := u.GetUsersRestaurants(user.ID)
		if err != nil {
			writeErrorResponse([]entities.APIError{
				{
					Message: err.Error(),
					Code:    "",
				},
			}, http.StatusBadRequest, w)
			return
		}

		writeSuccesResponse(rs, w)
		return
	}

}

type getReportsRequestBody struct {
	RestaurantID  uint   `json:"restaurant_id" valid:"optional"`
	StartInterval string `json:"start_interval" valid:"type(string),required"`
	EndInterval   string `json:"end_interval" valid:"type(string),required"`
}

func (h *ReportsHandler) GetRestaurantReport(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequestContext(r, w)
	if user.UserRole.Name != "admin" && user.UserRole.Name != "package_supplier" && user.UserRole.Name != "restaurant" && user.UserRole.Name != "closed_loop_restaurant" {
		writeErrorResponse([]entities.APIError{}, http.StatusForbidden, w)
		return
	}

	b := getReportsRequestBody{}
	if !validateRequestBody(r, w, &b) {
		return
	}

	u := usecases.ReportsUsecase{
		ReportsRepository:     repositories.ReportsRepository{DB: h.DB},
		RestaurantsRepository: repositories.RestaurantsRepository{DB: h.DB},
	}

	if user.UserRole.Name != "package_supplier" && user.UserRole.Name != "admin" && b.RestaurantID == 0 {
		writeErrorResponse([]entities.APIError{
			{
				Message: "Call only allowed for package suppliers",
				Code:    "1",
			},
		}, http.StatusForbidden, w)
		return
	}

	if user.UserRole.Name == "admin" && b.RestaurantID == 0 {
		rp, err := u.GetAdminTotalsReport(user.ID, b.StartInterval, b.EndInterval)
		if err != nil {
			writeErrorResponse([]entities.APIError{
				{
					Message: err.Error(),
					Code:    "2",
				},
			}, http.StatusBadRequest, w)
			return
		}

		writeSuccesResponse(rp, w)
		return
	}

	if user.UserRole.Name == "admin" && b.RestaurantID != 0 {
		rp, err := u.GetRestaurantReport(user.ID, b.RestaurantID, b.StartInterval, b.EndInterval)
		if err != nil {
			writeErrorResponse([]entities.APIError{
				{
					Message: err.Error(),
					Code:    "3",
				},
			}, http.StatusBadRequest, w)
			return
		}

		writeSuccesResponse(rp, w)
		return
	}

	if user.UserRole.Name == "package_supplier" && b.RestaurantID == 0 {
		rp, err := u.GetSupplierTotalsReport(user.ID, b.StartInterval, b.EndInterval)
		if err != nil {
			writeErrorResponse([]entities.APIError{
				{
					Message: err.Error(),
					Code:    "2",
				},
			}, http.StatusBadRequest, w)
			return
		}

		writeSuccesResponse(rp, w)
		return
	}

	if user.UserRole.Name == "package_supplier" && b.RestaurantID != 0 {
		rp, err := u.GetSupplierRestaurantReport(user.ID, b.RestaurantID, b.StartInterval, b.EndInterval)
		if err != nil {
			writeErrorResponse([]entities.APIError{
				{
					Message: err.Error(),
					Code:    "3",
				},
			}, http.StatusBadRequest, w)
			return
		}

		writeSuccesResponse(rp, w)
		return
	}

	rp, err := u.GetRestaurantReport(user.ID, b.RestaurantID, b.StartInterval, b.EndInterval)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "4",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(rp, w)
	return
}
