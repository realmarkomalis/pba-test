package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type ErrorResponse struct {
	Errors []ResponseError `json:"errors"`
}

type ResponseError struct {
	Message    string `json:"message"`
	Code       string `json:"code"`
	ResourceID uint   `json:"resource_id"`
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

func writeErrorResponse(errors []ResponseError, statusCode int, w http.ResponseWriter) {
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
		writeErrorResponse([]ResponseError{
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
		writeErrorResponse([]ResponseError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return false
	}

	result, err := govalidator.ValidateStruct(body)
	if err != nil {
		writeErrorResponse([]ResponseError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return false
	}

	if !result {
		writeErrorResponse([]ResponseError{
			{
				Message: err.Error(),
				Code:    "0",
			},
		}, http.StatusBadRequest, w)
		return false
	}

	return true
}

func addCookie(w http.ResponseWriter, name, value string, ttl time.Duration) {
	expire := time.Now().Add(ttl)
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expire,
		HttpOnly: true,
		Domain:   "127.0.0.1",
	}
	http.SetCookie(w, &cookie)
}
