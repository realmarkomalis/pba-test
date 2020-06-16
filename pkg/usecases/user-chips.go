package usecases

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

const ReturnEntryPoints = 6

type UserChipsUsecase struct {
	UserReturnRepo repositories.IUserReturnRepository
}

func (u UserChipsUsecase) GetUserChips(userID uint) (*entities.Chips, error) {
	entries, err := u.UserReturnRepo.GetUserReturnEntries(userID)
	if err != nil {
		return nil, err
	}

	total := 0
	for _, e := range entries {
		for _, r := range e.Returns {
			if r.Status == "Fulfilled" {
				total = total + ReturnEntryPoints
			}
		}
	}
	return &entities.Chips{
		User: entities.User{
			ID: userID,
		},
		Total: uint(total),
	}, nil
}
