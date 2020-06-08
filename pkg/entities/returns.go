package entities

import "time"

type Return struct {
	ID              uint            `json:"id"`
	Status          string          `json:"status"`
	Package         Package         `json:"package"`
	PickupRequest   PickupRequest   `json:"pickup_request"`
	PackageDispatch PackageDispatch `json:"package_dispatch"`
	PackagePickup   PackagePickup   `json:"package_pickup"`
	CreatedAt       time.Time       `json:"created_at"`
}

type PackageDispatch struct {
	ID        uint      `json:"id"`
	User      User      `json:"user"`
	ReturnID  uint      `json:"return_id"`
	CreatedAt time.Time `json:"created_at"`
}

type PickupRequest struct {
	ID          uint       `json:"id"`
	User        User       `json:"user"`
	ReturnID    uint       `json:"return_id"`
	PickupUser  User       `json:"pickup_user"`
	RequestUser User       `json:"request_user"`
	PickupSlot  PickupSlot `json:"pickup_slot"`
	CreatedAt   time.Time  `json:"created_at"`
}

type PackagePickup struct {
	ID        uint      `json:"id"`
	User      User      `json:"user"`
	ReturnID  uint      `json:"return_id"`
	CreatedAt time.Time `json:"created_at"`
}
