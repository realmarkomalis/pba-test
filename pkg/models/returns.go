package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ReturnStatus int

func (r ReturnStatus) String() string {
	return [...]string{"Created", "Dispatched", "Scheduled", "Fulfilled"}[r]
}

const (
	Created ReturnStatus = iota
	Dispatched
	Scheduled
	Fulfilled
)

type Return struct {
	gorm.Model
	Status          ReturnStatus
	Package         Package
	PackageID       uint
	PackageDispatch PackageDispatch
	ReturnRequest   ReturnRequest
	PackagePickup   PackagePickup
}

type UserReturnEntry struct {
	gorm.Model
	User    User
	UserID  uint
	Returns []Return `gorm:"many2many:user_return_entry_returns"`
}

type UserReturnEntryReturns struct {
	gorm.Model
	UserReturnEntryID uint
	ReturnID          uint
}

type PackageDispatch struct {
	gorm.Model
	User     User
	UserID   uint
	ReturnID uint
}

type ReturnRequest struct {
	gorm.Model
	User         User
	UserID       uint
	ReturnID     uint
	PickupSlot   PickupSlot
	PickupSlotID uint
}

type PackagePickup struct {
	gorm.Model
	User     User
	UserID   uint
	ReturnID uint
}

type PickupSlot struct {
	gorm.Model
	StartDateTime time.Time
	EndDateTime   time.Time
	User          User
	UserID        uint `gorm:"-"`
	Booked        bool `gorm:"default:false"`
}
