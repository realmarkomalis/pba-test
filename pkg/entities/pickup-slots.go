package entities

import "time"

type PickupSlot struct {
	ID              uint      `json:"id"`
	StartDateTime   time.Time `json:"start_datetime"`
	EndDateTime     time.Time `json:"end_datetime"`
	Booked          bool      `json:"booked"`
	ReturnRequestID uint      `json:"return_request_id"`
}
