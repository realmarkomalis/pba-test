package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type ErrorResponse struct {
	Errors []entities.APIError `json:"errors"`
}

func writeSuccesResponse(data interface{}, w http.ResponseWriter) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func writeErrorResponse(errors []entities.APIError, statusCode int, w http.ResponseWriter) {
	js, err := json.Marshal(ErrorResponse{Errors: errors})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getUserFromRequestContext(r *http.Request, w http.ResponseWriter) *entities.User {
	user := r.Context().Value("User")

	if user == nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: "No user found",
				Code:    "0",
			},
		}, http.StatusForbidden, w)
		return nil
	}

	return user.(*entities.User)
}

func validateRequestBody(r *http.Request, w http.ResponseWriter, body interface{}) bool {
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return false
	}

	result, err := govalidator.ValidateStruct(body)
	if err != nil {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return false
	}

	if !result {
		writeErrorResponse([]entities.APIError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return false
	}

	return true
}
