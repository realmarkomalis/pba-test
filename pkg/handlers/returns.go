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
	PackageIDs []uint `json:"package_ids" valid:"required"`
}

type createReturnRequestRequestBody struct {
	PackageIDs []uint `json:"package_ids" valid:"required"`
	SlotID     uint   `json:"slot_id" valid:"type(uint)"`
}

type createPackagePickupRequestBody struct {
	PackageIDs []uint `json:"package_ids" valid:"required"`
}

type returnsResponseBody struct {
	Returns []*entities.Return     `json:"returns"`
	Errors  map[uint]ErrorResponse `json:"errors"`
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
	if !validateRequestBody(r, w, &b) {
		return
	}

	user := getUserFromRequestContext(r, w)
	rr := repositories.ReturnRepository{DB: h.DB}
	urr := repositories.UserReturnRepository{DB: h.DB}
	pr := repositories.PickupSlotsRepository{DB: h.DB}
	u := usecases.ReturnsUsecase{
		ReturnsRepo:     rr,
		UserReturnRepo:  urr,
		PickupSlotsRepo: pr,
	}

	returns, err := u.CreatePackageDispatches(user.ID, b.PackageIDs)
	if err != nil {
		return
	}

	writeSuccesResponse(returnsResponseBody{Returns: returns}, w)
}

func (h ReturnsRequestsHandler) CreateReturnRequest(w http.ResponseWriter, r *http.Request) {
	b := createReturnRequestRequestBody{}
	if !validateRequestBody(r, w, &b) {
		return
	}

	user := getUserFromRequestContext(r, w)
	rr := repositories.ReturnRepository{DB: h.DB}
	urr := repositories.UserReturnRepository{DB: h.DB}
	pr := repositories.PickupSlotsRepository{DB: h.DB}
	u := usecases.ReturnsUsecase{
		ReturnsRepo:     rr,
		UserReturnRepo:  urr,
		PickupSlotsRepo: pr,
	}

	returns, err := u.CreateReturnRequests(user.ID, b.SlotID, b.PackageIDs)
	if err != nil {
		return
	}

	writeSuccesResponse(returnsResponseBody{Returns: returns}, w)
}

func (h ReturnsRequestsHandler) CreatePackagePickup(w http.ResponseWriter, r *http.Request) {
	b := createPackagePickupRequestBody{}
	if !validateRequestBody(r, w, &b) {
		return
	}

	user := getUserFromRequestContext(r, w)
	rr := repositories.ReturnRepository{DB: h.DB}
	urr := repositories.UserReturnRepository{DB: h.DB}
	pr := repositories.PickupSlotsRepository{DB: h.DB}
	u := usecases.ReturnsUsecase{
		ReturnsRepo:     rr,
		UserReturnRepo:  urr,
		PickupSlotsRepo: pr,
	}

	returns, err := u.CreatePackagePickups(user.ID, b.PackageIDs)
	if err != nil {
		return
	}

	writeSuccesResponse(returnsResponseBody{Returns: returns}, w)
}

func (h ReturnsRequestsHandler) GetReturns(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequestContext(r, w)
	rr := repositories.ReturnRepository{DB: h.DB}
	urr := repositories.UserReturnRepository{DB: h.DB}
	pr := repositories.PickupSlotsRepository{DB: h.DB}
	u := usecases.ReturnsUsecase{
		ReturnsRepo:     rr,
		UserReturnRepo:  urr,
		PickupSlotsRepo: pr,
	}
	rets, err := u.GetUserReturns(user.ID)
	if err != nil {

	}

	writeSuccesResponse(rets, w)
}

func (h ReturnsRequestsHandler) GetPackageDispatches(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequestContext(r, w)

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
	user := getUserFromRequestContext(r, w)

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
	user := getUserFromRequestContext(r, w)

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
