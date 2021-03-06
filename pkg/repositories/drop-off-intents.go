package repositories

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type DropOffIntentsRepository struct {
	DB *gorm.DB
}

func (r DropOffIntentsRepository) CreateDropOffIntent(userID, returnID, dropOffPointID uint) (*entities.DropOffIntent, error) {
	doi := models.DropOffIntent{
		UserID:         userID,
		ReturnID:       returnID,
		DropOffPointID: dropOffPointID,
	}

	result := r.DB
	if dropOffPointID > 0 {
		result = result.
			Create(&doi)
	} else {
		result = result.
			Select("id", "created_at", "updated_at", "user_id", "return_id").
			Create(&doi)
	}

	err := result.Error
	if err != nil {
		return nil, err
	}

	d := doi.ModelToEntity()
	return &d, nil
}

func (r DropOffIntentsRepository) GetUserDropOffIntents(userID uint) ([]entities.DropOffIntent, error) {
	doi := []models.DropOffIntent{}

	err := r.DB.
		Where("user_id = ?", userID).
		Find(&doi).
		Error
	if err != nil {
		return nil, err
	}

	dois := []entities.DropOffIntent{}
	for _, d := range doi {
		dois = append(dois, d.ModelToEntity())
	}

	return dois, nil
}
