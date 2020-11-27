package models

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type DropOffIntent struct {
	gorm.Model
	ReturnID       uint `gorm:"not null"`
	User           User
	UserID         uint `gorm:"not null"`
	DropOffPoint   DropOffPoint
	DropOffPointID uint `gorm:"default:null"`
}

func (d DropOffIntent) ModelToEntity() entities.DropOffIntent {
	return entities.DropOffIntent{
		ID:           d.ID,
		CreatedAt:    d.CreatedAt,
		ReturnID:     d.ReturnID,
		User:         d.User.ModelToEntity(),
		DropOffPoint: d.DropOffPoint.ModelToEntity(),
	}
}
