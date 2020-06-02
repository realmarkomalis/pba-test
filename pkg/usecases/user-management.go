package usecases

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type UserManagementUsecase struct {
}

func (u UserManagementUsecase) CreateNewUser(user entities.User) error {
	return nil
}
