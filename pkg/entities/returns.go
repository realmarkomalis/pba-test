package entities

import "time"

type Return struct {
	ID              uint            `json:"id"`
	Status          string          `json:"status"`
	StatusCode      int             `json:"status_code"`
	Package         Package         `json:"package"`
	PickupRequest   PickupRequest   `json:"pickup_request"`
	PackageDispatch PackageDispatch `json:"package_dispatch"`
	PackagePickup   PackagePickup   `json:"package_pickup"`
	CreatedAt       time.Time       `json:"created_at"`
}

type UserReturnEntry struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user"`
	Returns   []Return  `json:"returns"`
}

type PackageDispatch struct {
	ID        uint      `json:"id"`
	User      User      `json:"user"`
	ReturnID  uint      `json:"return_id"`
	CreatedAt time.Time `json:"created_at"`
}

type PickupRequest struct {
	ID          uint      `json:"id"`
	User        User      `json:"user"`
	ReturnID    uint      `json:"return_id"`
	PickupUser  User      `json:"pickup_user"`
	RequestUser User      `json:"request_user"`
	CreatedAt   time.Time `json:"created_at"`
}

type PackagePickup struct {
	ID        uint      `json:"id"`
	User      User      `json:"user"`
	ReturnID  uint      `json:"return_id"`
	CreatedAt time.Time `json:"created_at"`
}
