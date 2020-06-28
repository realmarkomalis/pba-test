package repositories

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type IUserReturnRepository interface {
	GetUserReturns(userID uint) ([]*entities.UserReturnEntry, error)
	GetUserReturnEntries(userID uint) ([]*entities.UserReturnEntry, error)
	GetScheduledReturnEntries() ([]*entities.UserReturnEntry, error)

	CreateUserReturnEntry(userID uint, returnIDs []uint) (*entities.UserReturnEntry, error)
}

type UserReturnRepository struct {
	DB *gorm.DB
}

func (r UserReturnRepository) CreateUserReturnEntry(userID uint, returnIDs []uint) (*entities.UserReturnEntry, error) {
	rets := []models.Return{}
	err := r.DB.
		Preload("Package").
		Where("id in (?)", returnIDs).
		Find(&rets).
		Error
	if err != nil {
		return nil, err
	}

	ure := models.UserReturnEntry{UserID: userID}
	err = r.DB.
		Save(&ure).
		Error
	if err != nil {
		return nil, err
	}

	for _, re := range rets {
		err = r.DB.
			Model(&ure).
			Association("Returns").
			Append(re).
			Error
		if err != nil {
			return nil, err
		}
	}

	err = r.DB.
		Save(&ure).
		Error
	if err != nil {
		return nil, err
	}

	return &entities.UserReturnEntry{
		ID:        ure.ID,
		CreatedAt: ure.CreatedAt,
		User: entities.User{
			ID: ure.User.ID,
		},
		Returns: returnModelsToEntities(rets),
	}, nil
}

func (r UserReturnRepository) GetUserReturns(userID uint) ([]*entities.UserReturnEntry, error) {
	return r.GetUserReturnEntries(userID)
}

func (r UserReturnRepository) GetUserReturnEntries(userID uint) ([]*entities.UserReturnEntry, error) {
	urs := []models.UserReturnEntry{}
	err := r.DB.
		Model(&urs).
		Preload("User").
		Preload("Returns").
		Preload("Returns.Package").
		Where("user_id = ?", userID).
		Find(&urs).
		Error
	if err != nil {
		return nil, err
	}

	userReturns := []*entities.UserReturnEntry{}
	for _, ur := range urs {
		userReturns = append(userReturns, &entities.UserReturnEntry{
			ID:        ur.ID,
			CreatedAt: ur.CreatedAt,
			Returns:   returnModelsToEntities(ur.Returns),
		})
	}

	return userReturns, nil
}

func (r UserReturnRepository) GetScheduledReturnEntries() ([]*entities.UserReturnEntry, error) {
	rr := []models.Return{}
	err := r.DB.
		Model(&rr).
		Where("status = ?", 2).
		Find(&rr).
		Error
	if err != nil {
		return nil, err
	}

	urs := []models.UserReturnEntry{}
	err = r.DB.
		Model(&rr).
		Preload("Returns").
		Preload("Returns.Package").
		Preload("User").
		Preload("User.UserAddresses").
		Preload("User.UserRole").
		Related(&urs, "UserReturnEntries").
		Select("DISTINCT(id)").
		Error
	if err != nil {
		return nil, err
	}

	userReturns := []*entities.UserReturnEntry{}
	for _, ur := range urs {
		for i, rt := range ur.Returns {
			rq := models.ReturnRequest{}
			err := r.DB.
				Model(&rq).
				Preload("PickupSlot").
				Preload("User").
				Preload("User.UserAddresses").
				Where("return_id = ?", rt.ID).
				First(&rq).
				Error
			if err == nil {
				ur.Returns[i].ReturnRequest = rq
			}
		}

		userReturns = append(userReturns, &entities.UserReturnEntry{
			ID:        ur.ID,
			CreatedAt: ur.CreatedAt,
			User:      ur.User.ModelToEntity(),
			Returns:   returnModelsToEntities(ur.Returns),
		})
	}

	return userReturns, nil
}

func returnModelsToEntities(rets []models.Return) []entities.Return {
	returns := []entities.Return{}
	for _, r := range rets {
		returns = append(returns, r.ModelToEntity())
	}

	return returns
}
