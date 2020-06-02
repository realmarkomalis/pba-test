package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

func writeSuccesResponse(data interface{}, w http.ResponseWriter) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getUserFromRequestContext(r *http.Request, w http.ResponseWriter) *entities.User {
	user := r.Context().Value("User")

	if user == nil {
		http.Error(w, "No user found", http.StatusInternalServerError)
		return nil
	}

	return user.(*entities.User)
}

func validateRequestBody(r *http.Request, w http.ResponseWriter, body interface{}) {
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := govalidator.ValidateStruct(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !result {
		http.Error(w, "No result reading the request body", http.StatusBadRequest)
		return
	}
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
