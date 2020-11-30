package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/usecases"
)

type DepositsHandler struct {
	DB *gorm.DB
}

func (h DepositsHandler) GetUserDepositItems(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	if u == nil {
		return
	}

	du := usecases.DepositsUsecase{
		ReturnsRepository:        repositories.ReturnRepository{DB: h.DB},
		DropOffIntentsRepository: repositories.DropOffIntentsRepository{DB: h.DB},
		DepositsRepository:       repositories.DepositsRepository{DB: h.DB},
	}

	deposits, err := du.GetUserDepositItems(u.ID)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(deposits, w)
}

func (h DepositsHandler) ListCustomerDepositPayouts(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	if u == nil {
		return
	}

	if u.UserRole.Name != "admin" {
		writeErrorResponse([]entities.APIError{}, http.StatusForbidden, w)
		return
	}

	du := usecases.DepositsUsecase{
		ReturnsRepository:        repositories.ReturnRepository{DB: h.DB},
		DropOffIntentsRepository: repositories.DropOffIntentsRepository{DB: h.DB},
		DepositsRepository:       repositories.DepositsRepository{DB: h.DB},
	}

	ps, err := du.ListUnpaidCustomerDepositPayouts()
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(ps, w)
}

func (h DepositsHandler) CreateCustomerDepositPayout(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	if u == nil {
		return
	}

	if u.UserRole.Name != "customer" {
		writeErrorResponse([]entities.APIError{}, http.StatusForbidden, w)
		return
	}

	du := usecases.DepositsUsecase{
		ReturnsRepository:        repositories.ReturnRepository{DB: h.DB},
		DropOffIntentsRepository: repositories.DropOffIntentsRepository{DB: h.DB},
		DepositsRepository:       repositories.DepositsRepository{DB: h.DB},
	}

	ps, err := du.ClaimCustomerDeposits(u.ID)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(ps, w)
}

type updateCustomerDepositPayoutRequestBody struct {
	ReceivingUserID uint    `json:"receiving_user_id" valid:"required"`
	PaymentID       string  `json:"payment_id" valid:"required"`
	Amount          float64 `json:"amount" valid:"required"`
}

func (h DepositsHandler) UpdateCustomerDepositPayout(w http.ResponseWriter, r *http.Request) {
	u := getUserFromRequestContext(r, w)
	if u == nil {
		return
	}

	if u.UserRole.Name != "admin" {
		writeErrorResponse([]entities.APIError{}, http.StatusForbidden, w)
		return
	}

	b := updateCustomerDepositPayoutRequestBody{}
	if !validateRequestBody(r, w, &b) {
		return
	}

	du := usecases.DepositsUsecase{
		ReturnsRepository:        repositories.ReturnRepository{DB: h.DB},
		DropOffIntentsRepository: repositories.DropOffIntentsRepository{DB: h.DB},
		DepositsRepository:       repositories.DepositsRepository{DB: h.DB},
	}

	ps, err := du.PayCustomerDeposits(b.ReceivingUserID, u.ID, b.PaymentID)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return
	}

	writeSuccesResponse(ps, w)
}
