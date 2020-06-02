package usecases

import (
	"fmt"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
	"gitlab.com/markomalis/packback-api/pkg/services"
)

type AuthUsecase struct {
	UserRepo              repositories.IUserRepository
	AuthService           services.AuthService
	UserManagementService services.UserManagementService
	TokenService          services.TokenService
}

func (a AuthUsecase) RequestLogin(email string) (bool, error) {
	user, err := a.UserRepo.GetByEmail(email)
	if err != nil {
		fmt.Printf("No user found with email %s, creating user in IDP \n", email)
		_, err = a.UserManagementService.CreateUser(email)
		if err != nil {
			fmt.Printf("Error creating user in IDP %s \n", err)
			return false, err
		}

		user, err = a.UserRepo.Create(&entities.User{Email: email})
		if err != nil {
			fmt.Printf("Error creating user in DB: %s \n", err)
			return false, err
		}
	}

	fmt.Printf("Requesting login for user %s \n", user.Email)
	err = a.AuthService.RequestLogin(user.Email)
	if err != nil {
		fmt.Printf("Error requesting login %s \n", err)
		return false, err
	}

	return true, nil
}

func (a AuthUsecase) LoginUserWithCode(email, code string) (*entities.User, error) {
	user, err := a.AuthService.Login(email, code)
	if err != nil {
		fmt.Printf("Error loging in user with email %s: %s \n", email, err)
		return nil, err
	}

	repoUser, err := a.UserRepo.GetByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	token, err := a.TokenService.Encode(repoUser)
	if err != nil {
		return nil, err
	}

	repoUser.Token = token

	return repoUser, nil
}
