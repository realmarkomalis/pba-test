package models

import (
	"github.com/jinzhu/gorm"
)

type Supplier struct {
	gorm.Model
	Name         string
	User         User
	UserID       uint
	PackageTypes []PackageType
}
