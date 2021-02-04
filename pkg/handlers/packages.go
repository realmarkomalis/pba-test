package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/services"
)

type PackagesHandler struct {
	DB *gorm.DB
}

type createOrUpdatePackageBody struct {
	PackageCode   string `json:"package_code"`
	PackageTypeID uint   `json:"package_type_id"`
}

func (h PackagesHandler) CreateOrUpdatePackage(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	if u == nil {
		return
	}

	if u.UserRole.Name != "admin" {
		writeErrorResponse([]entities.APIError{}, http.StatusForbidden, w)
		return
	}

	b := createOrUpdatePackageBody{}
	if !validateRequestBody(r, w, &b) {
		return
	}

	pr := repositories.PackagesRepository{DB: h.DB}
	p, err := pr.CreateOrUpdatePackage(b.PackageCode, b.PackageTypeID)
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

func (h PackagesHandler) GetPackage(w http.ResponseWriter, r *http.Request) {
	packageID := mux.Vars(r)["package_id"]
	u, err := uuid.FromString(packageID)
	if err != nil {
		h.GetPackageByID(w, packageID)
	} else {
		h.GetPackageByCode(w, u.String())
	}
}

func (h PackagesHandler) GetPackageByID(w http.ResponseWriter, packageID string) {
	id, err := strconv.ParseUint(packageID, 10, 64)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	pr := repositories.PackagesRepository{DB: h.DB}
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

func (h PackagesHandler) GetPackageByCode(w http.ResponseWriter, code string) {
	pr := repositories.PackagesRepository{DB: h.DB}
	p, err := pr.GetPackageByCode(code)
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

type getPackageCodeResponse struct {
	Code string `json:"code"`
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
