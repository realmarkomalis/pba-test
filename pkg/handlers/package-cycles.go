package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type PackageCyclesHandler struct {
	DB *gorm.DB
}

func (h PackageCyclesHandler) GetPackageCycles(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	if u == nil {
		return
	}

	if u.UserRole.Name == "admin" {
		writeErrorResponse([]entities.APIError{
			{
				Message: errors.New("Laters mattie!").Error(),
				Code:    "0",
			},
		}, http.StatusForbidden, w)
		return
	}

	pr := repositories.PackagesRepository{DB: h.DB}
	rr := repositories.ReturnRepository{DB: h.DB}

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

	rets, err := rr.GetReturns(p.ID)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusNotFound, w)
		return
	}

	writeSuccesResponse(rets, w)
}
