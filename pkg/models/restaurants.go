package models

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type Restaurant struct {
	gorm.Model
	Name   string
	User   User
	UserID uint
}

func (r Restaurant) ModelToEntity() entities.Restaurant {
	return entities.Restaurant{
		ID:   r.ID,
		Name: r.Name,
		User: r.User.ModelToEntity(),
	}
}
