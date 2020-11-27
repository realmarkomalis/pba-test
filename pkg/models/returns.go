package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type ReturnStatus int

func (r ReturnStatus) String() string {
	return [...]string{"Created", "Dispatched", "Scheduled", "Fulfilled", "Collected", "DropOffIntended"}[r]
}

const (
	Created ReturnStatus = iota
	Dispatched
	Scheduled
	Fulfilled
	Collected
	DropOffIntended
)

type Return struct {
	gorm.Model
	Status            ReturnStatus
	Package           Package
	PackageID         uint
	PackageSupply     PackageSupply
	PackageDispatch   PackageDispatch
	ReturnRequest     ReturnRequest
	PackagePickup     PackagePickup
	UserReturnEntries []UserReturnEntry `gorm:"many2many:user_return_entry_returns"`
}

func (r Return) ModelToEntity() entities.Return {
	return entities.Return{
		ID:              r.ID,
		CreatedAt:       r.CreatedAt,
		Status:          r.Status.String(),
		StatusCode:      int(r.Status),
		Package:         r.Package.ModelToEntity(),
		PackageSupply:   r.PackageSupply.ModelToEntity(),
		PackageDispatch: r.PackageDispatch.ModelToEntity(),
		PickupRequest:   r.ReturnRequest.ModelToEntity(),
		PackagePickup:   r.PackagePickup.ModelToEntity(),
	}
}

type UserReturnEntry struct {
	gorm.Model
	User    User
	UserID  uint
	Returns []Return `gorm:"many2many:user_return_entry_returns"`
}

func (r UserReturnEntry) ModelToEntity() entities.UserReturnEntry {
	rs := []entities.Return{}
	for _, rr := range r.Returns {
		rs = append(rs, rr.ModelToEntity())
	}

	return entities.UserReturnEntry{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
		User:      r.User.ModelToEntity(),
		Returns:   rs,
	}
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

func (pd PackageDispatch) PackageDispatch() entities.PackageDispatch {
	return entities.PackageDispatch{
		ID:        pd.ID,
		User:      pd.User.ModelToEntity(),
		ReturnID:  pd.ReturnID,
		CreatedAt: pd.CreatedAt,
	}
}

func (pd PackageDispatch) ModelToEntity() entities.PackageDispatch {
	return entities.PackageDispatch{
		ID:        pd.ID,
		User:      pd.User.ModelToEntity(),
		ReturnID:  pd.ReturnID,
		CreatedAt: pd.CreatedAt,
	}
}

type ReturnRequest struct {
	gorm.Model
	User         User
	UserID       uint
	ReturnID     uint
	PickupSlot   PickupSlot
	PickupSlotID uint
	PickupNote   string
}

func (rr ReturnRequest) ModelToEntity() entities.PickupRequest {
	return entities.PickupRequest{
		ID:         rr.ID,
		User:       rr.User.ModelToEntity(),
		ReturnID:   rr.ReturnID,
		PickupSlot: rr.PickupSlot.ModelToEntity(),
		PickupNote: rr.PickupNote,
		CreatedAt:  rr.CreatedAt,
	}
}

type PackagePickup struct {
	gorm.Model
	User     User
	UserID   uint
	ReturnID uint
}

func (pp PackagePickup) ModelToEntity() entities.PackagePickup {
	return entities.PackagePickup{
		ID:        pp.ID,
		User:      pp.User.ModelToEntity(),
		ReturnID:  pp.ReturnID,
		CreatedAt: pp.CreatedAt,
	}
}

type PickupSlot struct {
	gorm.Model
	StartDateTime time.Time
	EndDateTime   time.Time
	User          User
	UserID        uint `gorm:"-"`
	Booked        bool `gorm:"default:false"`
}

func (ps PickupSlot) ModelToEntity() entities.PickupSlot {
	return entities.PickupSlot{
		ID:            ps.ID,
		StartDateTime: ps.StartDateTime,
		EndDateTime:   ps.EndDateTime,
		Booked:        ps.Booked,
	}
}

type PackageSupply struct {
	gorm.Model
	User         User
	UserID       uint
	ReturnID     uint
	Restaurant   Restaurant
	RestaurantID uint
}

func (ps PackageSupply) ModelToEntity() entities.PackageSupply {
	return entities.PackageSupply{
		ID:         ps.ID,
		User:       ps.User.ModelToEntity(),
		ReturnID:   ps.ReturnID,
		CreatedAt:  ps.CreatedAt,
		Restaurant: ps.Restaurant.ModelToEntity(),
	}
}
