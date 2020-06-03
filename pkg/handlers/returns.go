package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/usecases"
)

type ReturnsRequestsHandler struct {
	DB *gorm.DB
}

type createReturnRequestBody struct {
	PackageID uint `json:"package_id" valid:"type(uint)"`
}

type createReturnRequestRequestBody struct {
	SlotID uint `json:"slot_id" valid:"type(uint)"`
}

type createPickupRequestBody struct {
}

func (h ReturnsRequestsHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(*entities.User)
	if user == nil {
		http.Error(w, "No user found", http.StatusInternalServerError)
		return
	}

	rr := repositories.ReturnRepository{h.DB}
	u := usecases.ReturnsUsecase{rr}

	body := createReturnRequestBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := govalidator.ValidateStruct(body)
	if err != nil || !result {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ret, err := u.CreateNewReturn(user.ID, body.PackageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func (h ReturnsRequestsHandler) CreateReturnRequest(w http.ResponseWriter, r *http.Request) {
	b := createReturnRequestRequestBody{}
	validateRequestBody(r, w, &b)

	vars := mux.Vars(r)
	packageID, err := strconv.ParseUint(vars["package_id"], 10, 32)
	if err != nil {
		fmt.Println(err)
	}

	user := getUserFromRequestContext(r, w)
	pr := repositories.ReturnRepository{h.DB}
	ret, err := pr.CreateReturnRequest(user.ID, uint(packageID), b.SlotID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeSuccesResponse(ret, w)
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

	if user.UserRole.Name == "rider" {
		h.GetPackagePickups(w, r)
		return
	}

	if user.UserRole.Name == "customer" {
		h.GetReturnRequests(w, r)
		return
	}
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

func (h ReturnsRequestsHandler) CreatePackagePickup(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(*entities.User)
	if user == nil {
		http.Error(w, "No user found", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	packageID, err := strconv.ParseUint(vars["package_id"], 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	rr := repositories.ReturnRepository{h.DB}
	u := usecases.ReturnsUsecase{rr}

	ret, err := u.CreatePackagePickup(user.ID, uint(packageID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
