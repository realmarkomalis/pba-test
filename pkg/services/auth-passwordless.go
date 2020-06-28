package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	AccesToken string `json:"access_token"`
	IDToken    string `json:"id_token"`
}

type userInfoResponse struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type errorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (p PasswordlessAuth) Login(username, code string) (string, error) {
	url := fmt.Sprintf("https://%s/oauth/token", p.Domain)
	requestBody, err := json.Marshal(map[string]interface{}{
		"client_id":     p.ClientID,
		"client_secret": p.ClientSecret,
		"username":      username,
		"otp":           code,
		"grant_type":    "http://auth0.com/oauth/grant-type/passwordless/otp",
		"realm":         "email",
		"scope":         "openid profile email",
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Add("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", errors.New("Error logining in user")
	}

	body := loginResponse{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", err
	}

	return body.AccesToken, nil
}

func (p PasswordlessAuth) UserInfo(token string) (*entities.User, error) {
	url := fmt.Sprintf("https://%s/userinfo", p.Domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, errors.New("Error logining in user")
	}

	body := userInfoResponse{}
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
		"send":          "code",
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

func decodeErrorResponse(b io.ReadCloser) {
	bodyBytes, err := ioutil.ReadAll(b)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}
