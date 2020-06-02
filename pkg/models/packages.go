package models

import (
	"github.com/jinzhu/gorm"
)

type Package struct {
	gorm.Model
	Name string
}
