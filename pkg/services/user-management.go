package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type UserManagementService struct {
	Domain       string
	ClientID     string
	ClientSecret string
	Audience     string
	AccessToken  string
}

type requestAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type createUserResponse struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Blocked       bool   `json:"blocked"`
}

func (u UserManagementService) requestAccessToken() (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"grant_type":    "client_credentials",
		"client_id":     u.ClientID,
		"client_secret": u.ClientSecret,
		"audience":      u.Audience,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s/oauth/token", u.Domain)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", errors.New("Unable to get access token")
	}

	body := requestAccessTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", err
	}

	return body.AccessToken, nil
}

func (u UserManagementService) CreateUser(email string) (*entities.User, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"email":          email,
		"email_verified": true,
		"verify_email":   false,
		"connection":     "email",
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s/api/v2/users", u.Domain)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	at, err := u.requestAccessToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", at))
	req.Header.Add("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, errors.New("Unable to create user")
	}

	body := createUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return &entities.User{Email: body.Email}, nil
}

func (u UserManagementService) GetUser(email string) (string, error) {
	return "", nil
}
