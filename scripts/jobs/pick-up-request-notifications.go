package jobs

import (
	"fmt"
	"time"

	"gitlab.com/markomalis/packback-api/pkg/models"
)

type PickupRequestNotificationRow struct {
	Email              string
	StreetName         string
	StreetNumber       string
	StreetNumberSuffix string
	PostalCode         string
	City               string
	SlotStart          time.Time
	SlotEnd            time.Time
}

func SendPickupRequestNotifications(j *Job) {
	fmt.Println("SendDropOffIntentNotifications")
	now := time.Now().UTC()

	rows, err := j.DB.
		Model(models.User{}).
		Select("users.email, user_addresses.street_name, user_addresses.house_number, user_addresses.house_number_suffix, user_addresses.postal_code, user_addresses.city, pickup_slots.start_date_time, pickup_slots.end_date_time").
		Joins("JOIN user_addresses on users.id = user_addresses.user_id").
		Joins("JOIN return_requests on return_requests.user_id = users.id").
		Joins("JOIN pickup_slots on pickup_slots.id = return_requests.pickup_slot_id").
		Where("return_requests.created_at BETWEEN ? AND ?", j.LowerBound, now).
		Group("users.email, user_addresses.street_name, user_addresses.house_number, user_addresses.house_number_suffix, user_addresses.postal_code, user_addresses.city, pickup_slots.start_date_time, pickup_slots.end_date_time").
		Rows()
	defer rows.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		row := PickupRequestNotificationRow{}
		rows.Scan(&row.Email, &row.StreetName, &row.StreetNumber, &row.StreetNumberSuffix, &row.PostalCode, &row.City, &row.SlotStart, &row.SlotEnd)
		fmt.Println("SCAN ROW")
		fmt.Println(row)
		sendPickupRequestMail(j, &row)
	}

	j.SetLowerBound(now)
}

func sendPickupRequestMail(j *Job, row *PickupRequestNotificationRow) error {
	fmt.Printf("SendDropOffIntentNotifications mail user with email: %s\n", row.Email)
	address := fmt.Sprintf("%s %s%s", row.StreetName, row.StreetNumber, row.StreetNumberSuffix)
	start, end := row.SlotStart, row.SlotEnd
	date := fmt.Sprintf("%02d/%02d/%d %02d:%02d - %02d:%02d", start.Day(), start.Month(), start.Year(), start.Hour(), start.Minute(), end.Hour(), end.Minute())

	j.MailService.SendTemplateEmail(
		"\"PackBack\" <info@packback.app>",
		[]string{row.Email},
		"Lekker bezig!",
		"return-requested",
		map[string]interface{}{
			"is_pickup":   true,
			"address":     address,
			"city":        row.City,
			"postal_code": row.PostalCode,
			"date":        date,
		},
	)

	return nil
}
