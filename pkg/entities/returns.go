package entities

import "time"

type Return struct {
	ID      uint    `json:"id"`
	Status  string  `json:"status"`
	Package Package `json:"package"`
	// PackageDispatch PackageDispatch `json:"package_dispatch"`
	// PickupRequest   PickupRequest   `json:"pickup_request"`
	// PackagePickup   PackagePickup   `json:"package_pickup"`
	CreatedAt time.Time `json:"created_at"`
}

type PackageDispatch struct {
	ID        uint      `json:"id"`
	User      User      `json:"user"`
	Return    Return    `json:"return"`
	CreatedAt time.Time `json:"created_at"`
}

type PickupRequest struct {
	ID          uint      `json:"id"`
	User        User      `json:"user"`
	PickupUser  User      `json:"pickup_user"`
	RequestUser User      `json:"request_user"`
	Return      Return    `json:"return"`
	CreatedAt   time.Time `json:"created_at"`
}

type PackagePickup struct {
	ID        uint      `json:"id"`
	User      User      `json:"user"`
	Return    Return    `json:"return"`
	CreatedAt time.Time `json:"created_at"`
}
