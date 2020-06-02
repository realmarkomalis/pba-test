package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type PasswordlessAuth struct {
	Domain       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type loginResponse struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type errorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (p PasswordlessAuth) Login(username, code string) (*entities.User, error) {
	url := fmt.Sprintf("https://%s/userinfo", p.Domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", code))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, errors.New("Error logining in user")
	}

	body := loginResponse{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return &entities.User{
		Email: body.Email,
	}, nil
}

func (p PasswordlessAuth) Logout(username string) error {
	return nil
}

func (p PasswordlessAuth) RequestLogin(username string) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"email":         username,
		"client_id":     p.ClientID,
		"client_secret": p.ClientSecret,
		"connection":    "email",
		"send":          "link",
		"authParams": map[string]interface{}{
			"scope":        "openid profile email",
			"redirect_uri": p.RedirectURL,
		},
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/passwordless/start", p.Domain)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("auth0-forwarded-for", "127.0.0.1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return errors.New("Error requesting login")
	}

	return nil
}

func decodeErrorResponse(b io.ReadCloser) errorResponse {
	body := errorResponse{}
	err := json.NewDecoder(b).Decode(&body)
	if err != nil {
		fmt.Println("Error decoding error response")
	}

	fmt.Printf("Error response: %s ; %s", body.Error, body.ErrorDescription)
	return body
}
