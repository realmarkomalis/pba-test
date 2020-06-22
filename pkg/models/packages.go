package models

import (
	"github.com/jinzhu/gorm"
)

type Package struct {
	gorm.Model
	Name string
	// PackageCode uuid.UUID `gorm:"type:uuid;not null; default:uuid_generate_v4()"`
}
