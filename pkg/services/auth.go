package services

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type AuthService interface {
	Login(username, password string) (*entities.User, error)
	Logout(username string) error
	RequestLogin(username string) error
}
