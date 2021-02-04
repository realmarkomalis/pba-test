package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type PackageTypesHandler struct {
	DB *gorm.DB
}

func (h PackageTypesHandler) GetPackageTypes(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequestContext(r, w)
	if user.UserRole.Name != "admin" {
		writeErrorResponse([]entities.APIError{}, http.StatusForbidden, w)
		return
	}

	ptr := repositories.PackageTypesRepository{DB: h.DB}
	ts, err := ptr.GetPackageTypes()
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusNotFound, w)
		return
	}

	writeSuccesResponse(ts, w)
}
