package repositories

import (
	"time"

	"github.com/jinzhu/gorm"

	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type PickupSlotsRepository struct {
	DB *gorm.DB
}

func (r PickupSlotsRepository) GetPickupSlot(slotID uint) (*entities.PickupSlot, error) {
	p := models.PickupSlot{}
	err := r.DB.
		Where("id = ?", slotID).
		First(&p).
		Error

	if err != nil {
		return nil, err
	}

	return &entities.PickupSlot{
		ID:            p.ID,
		StartDateTime: p.StartDateTime,
		EndDateTime:   p.EndDateTime,
		Booked:        p.Booked,
	}, nil
}

func (r PickupSlotsRepository) GetPickupSlots() ([]*entities.PickupSlot, error) {
	ps := []models.PickupSlot{}
	err := r.DB.
		Order("start_date_time asc").
		Where("start_date_time > ?", time.Now()).
		// Where("booked = ?", false).
		Find(&ps).
		Error

	if err != nil {
		return nil, err
	}

	slots := []*entities.PickupSlot{}
	for _, p := range ps {
		slots = append(slots, &entities.PickupSlot{
			ID:            p.ID,
			StartDateTime: p.StartDateTime,
			EndDateTime:   p.EndDateTime,
			Booked:        p.Booked,
		})
	}

	return slots, nil
}

func (r PickupSlotsRepository) BookPickupSlot(slotID uint) (*entities.PickupSlot, error) {
	p := models.PickupSlot{}
	err := r.DB.
		Model(&p).
		Where("id = ?", slotID).
		Update("booked", "true").
		Error
	if err != nil {
		return nil, err
	}

	return &entities.PickupSlot{
		ID:            p.ID,
		StartDateTime: p.StartDateTime,
		EndDateTime:   p.EndDateTime,
		Booked:        p.Booked,
	}, nil
}
