package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/services"
)

type PackagesHandler struct {
	DB *gorm.DB
}

type getPackageCodeResponse struct {
	Code string `json:"code"`
}

func (h PackagesHandler) GetPackage(w http.ResponseWriter, r *http.Request) {
	pr := repositories.PackagesRepository{DB: h.DB}

	id, err := strconv.ParseUint(mux.Vars(r)["package_id"], 10, 64)
	p, err := pr.GetPackage(uint(id))
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusNotFound, w)
		return
	}

	writeSuccesResponse(p, w)
}

func (h PackagesHandler) GetPackageCode(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	if u == nil {
		return
	}

	id, err := strconv.ParseUint(mux.Vars(r)["package_id"], 10, 64)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	s := services.QRCodeService{}
	code, err := s.CreateBase64Code(fmt.Sprintf("https://packback.app/scanner?packback_id=%d", id))
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusInternalServerError, w)
		return
	}

	writeSuccesResponse(getPackageCodeResponse{Code: code}, w)
}
