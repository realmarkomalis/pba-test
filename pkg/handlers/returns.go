package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/usecases"
)

type ReturnsRequestsHandler struct {
	DB *gorm.DB
}

type createPackageDispatchRequestBody struct {
	PackageIDs []uint `json:"package_ids"`
}

type createReturnRequestRequestBody struct {
	PackageIDs []uint `json:"package_ids"`
	SlotID     uint   `json:"slot_id" valid:"type(uint)"`
}

type createPackagePickupRequestBody struct {
	PackageIDs []uint `json:"package_ids"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type returnsResponseBody struct {
	Returns map[uint]*entities.Return `json:"returns"`
	Errors  map[uint]ErrorResponse    `json:"errors"`
}

func (h ReturnsRequestsHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequestContext(r, w)

	if user.UserRole.Name == "restaurant" {
		h.CreatePackageDispatch(w, r)
		return
	}

	if user.UserRole.Name == "customer" {
		h.CreateReturnRequest(w, r)
		return
	}

	if user.UserRole.Name == "rider" {
		h.CreatePackagePickup(w, r)
		return
	}
}

func (h ReturnsRequestsHandler) CreatePackageDispatch(w http.ResponseWriter, r *http.Request) {
	b := createPackageDispatchRequestBody{}
	validateRequestBody(r, w, &b)

	user := getUserFromRequestContext(r, w)
	rr := repositories.ReturnRepository{h.DB}
	u := usecases.ReturnsUsecase{rr}

	returns := make(map[uint]*entities.Return)
	errors := make(map[uint]ErrorResponse)
	for _, id := range b.PackageIDs {
		ret, err := u.CreatePackageDispatch(user.ID, id)
		if err != nil {
			errors[id] = ErrorResponse{
				Message: err.Error(),
				Code:    err.Error(),
			}
		} else {
			returns[id] = ret
		}
	}

	writeSuccesResponse(
		returnsResponseBody{
			Returns: returns,
			Errors:  errors,
		},
		w,
	)
}

func (h ReturnsRequestsHandler) CreateReturnRequest(w http.ResponseWriter, r *http.Request) {
	b := createReturnRequestRequestBody{}
	validateRequestBody(r, w, &b)

	user := getUserFromRequestContext(r, w)
	rr := repositories.ReturnRepository{h.DB}
	u := usecases.ReturnsUsecase{rr}

	returns := make(map[uint]*entities.Return)
	errors := make(map[uint]ErrorResponse)
	for _, id := range b.PackageIDs {
		ret, err := u.CreateReturnRequest(user.ID, id, b.SlotID)
		if err != nil {
			errors[id] = ErrorResponse{
				Message: err.Error(),
				Code:    err.Error(),
			}
		} else {
			returns[id] = ret
		}
	}

	writeSuccesResponse(
		returnsResponseBody{
			Returns: returns,
			Errors:  errors,
		},
		w,
	)
}

func (h ReturnsRequestsHandler) CreatePackagePickup(w http.ResponseWriter, r *http.Request) {
	b := createPackagePickupRequestBody{}
	validateRequestBody(r, w, &b)

	user := getUserFromRequestContext(r, w)
	rr := repositories.ReturnRepository{h.DB}
	u := usecases.ReturnsUsecase{rr}

	returns := make(map[uint]*entities.Return)
	errors := make(map[uint]ErrorResponse)
	for _, id := range b.PackageIDs {
		ret, err := u.CreatePackagePickup(user.ID, id)
		if err != nil {
			errors[id] = ErrorResponse{
				Message: err.Error(),
				Code:    err.Error(),
			}
		} else {
			returns[id] = ret
		}
	}

	writeSuccesResponse(
		returnsResponseBody{
			Returns: returns,
			Errors:  errors,
		},
		w,
	)
}

func (h ReturnsRequestsHandler) GetReturns(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(*entities.User)
	if user == nil {
		http.Error(w, "No user found", http.StatusInternalServerError)
		return
	}

	if user.UserRole.Name == "restaurant" {
		h.GetPackageDispatches(w, r)
		return
	}

	if user.UserRole.Name == "customer" {
		h.GetReturnRequests(w, r)
		return
	}

	if user.UserRole.Name == "rider" {
		h.GetPackagePickups(w, r)
		return
	}
}

func (h ReturnsRequestsHandler) GetPackageDispatches(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(*entities.User)
	if user == nil {
		http.Error(w, "No user found", http.StatusInternalServerError)
		return
	}

	rr := repositories.ReturnRepository{h.DB}
	returns, err := rr.GetPackageDispatches(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(returns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h ReturnsRequestsHandler) GetReturnRequests(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(*entities.User)
	if user == nil {
		http.Error(w, "No user found", http.StatusInternalServerError)
		return
	}

	rr := repositories.ReturnRepository{h.DB}
	returns, err := rr.GetReturnRequests(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(returns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h ReturnsRequestsHandler) GetPackagePickups(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(*entities.User)
	if user == nil {
		http.Error(w, "No user found", http.StatusInternalServerError)
		return
	}

	rr := repositories.ReturnRepository{h.DB}
	returns, err := rr.GetPackagePickups(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(returns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
