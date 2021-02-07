package jobs

import (
	"fmt"
	"time"

	"gitlab.com/markomalis/packback-api/pkg/models"
)

type DropOffIntentNotificationRow struct {
	Email              string
	DOPName            string
	StreetName         string
	StreetNumber       string
	StreetNumberSuffix string
	PostalCode         string
	City               string
}

func SendDropOffIntentNotifications(j *Job) {
	fmt.Println("SendDropOffIntentNotifications")
	now := time.Now().UTC()

	rows, err := j.DB.
		Model(models.User{}).
		Select("users.email, drop_off_points.name, drop_off_points.street_name, drop_off_points.house_number, drop_off_points.house_number_suffix, drop_off_points.postal_code, drop_off_points.city").
		Joins("JOIN drop_off_intents on users.id = drop_off_intents.user_id").
		Joins("JOIN drop_off_points on drop_off_points.id = drop_off_intents.drop_off_point_id").
		Where("drop_off_intents.created_at BETWEEN ? AND ?", j.LowerBound, now).
		Group("users.email, drop_off_points.name, drop_off_points.street_name, drop_off_points.house_number, drop_off_points.house_number_suffix, drop_off_points.postal_code, drop_off_points.city").
		Rows()
	defer rows.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		row := DropOffIntentNotificationRow{}
		rows.Scan(&row.Email, &row.DOPName, &row.StreetName, &row.StreetNumber, &row.StreetNumberSuffix, &row.PostalCode, &row.City)
		sendDropOffIntentMail(j, &row)
	}

	j.SetLowerBound(now)
}

func sendDropOffIntentMail(j *Job, row *DropOffIntentNotificationRow) error {
	fmt.Printf("SendDropOffIntentNotifications mail user with email: %s\n", row.Email)
	address := fmt.Sprintf("%s %s%s", row.StreetName, row.StreetNumber, row.StreetNumberSuffix)

	j.MailService.SendTemplateEmail(
		"\"PackBack\" <info@packback.app>",
		[]string{row.Email},
		"Lekker bezig!",
		"return-requested",
		map[string]interface{}{
			"is_drop_off": true,
			"name":        row.DOPName,
			"address":     address,
			"city":        row.City,
			"postal_code": row.PostalCode,
		},
	)

	return nil
}
