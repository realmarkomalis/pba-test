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
		Preload("ReturnRequest").
		Preload("ReturnRequest.PickupSlot").
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
		// Preload("Returns.PackageDispatch").
		// Preload("Returns.ReturnRequest").
		// Preload("Returns.ReturnRequest.PickupSlot").
		Preload("User").
		Preload("User.UserAddresses").
		Preload("User.UserRole").
		Related(&urs, "UserReturnEntries").
		Error
	if err != nil {
		return nil, err
	}

	userReturns := []*entities.UserReturnEntry{}
	for _, ur := range urs {
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

func returnModelToEntity(ret *models.Return) entities.Return {
	return entities.Return{
		ID:         ret.ID,
		CreatedAt:  ret.CreatedAt,
		Status:     ret.Status.String(),
		StatusCode: int(ret.Status),
		Package: entities.Package{
			ID:          ret.Package.ID,
			Name:        ret.Package.Name,
			PackageCode: ret.Package.PackageCode,
		},
		PickupRequest: entities.PickupRequest{
			ID: ret.ReturnRequest.ID,
			PickupSlot: entities.PickupSlot{
				ID:            ret.ReturnRequest.PickupSlot.ID,
				StartDateTime: ret.ReturnRequest.PickupSlot.StartDateTime,
				EndDateTime:   ret.ReturnRequest.PickupSlot.EndDateTime,
			},
			User: entities.User{
				ID: ret.ReturnRequest.User.ID,
			},
		},
	}
}
