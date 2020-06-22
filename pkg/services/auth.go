package services

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type AuthService interface {
	Login(username, password string) (string, error)
	Logout(username string) error
	RequestLogin(username string) error
	UserInfo(token string) (*entities.User, error)
}
