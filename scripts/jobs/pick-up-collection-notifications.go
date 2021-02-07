package jobs

import (
	"fmt"
	"time"

	"gitlab.com/markomalis/packback-api/pkg/models"
)

func SendPickUpCollectionNotifications(j *Job) {
	fmt.Println("SendPickUpCollectionNotifications")
	now := time.Now().UTC()
	users := []models.User{}

	err := j.DB.
		Select("users.email").
		Joins("JOIN return_requests on users.id = return_requests.user_id").
		Joins("JOIN package_pickups on package_pickups.return_id = return_requests.return_id").
		Where("package_pickups.created_at BETWEEN ? AND ?", j.LowerBound, now).
		Group("users.email").
		Find(&users).
		Error
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, u := range users {
		fmt.Printf("SendPickUpCollectionNotifications mail user with email: %s/n", u.Email)
		j.MailService.SendTemplateEmail(
			"\"PackBack\" <info@packback.app>",
			[]string{u.Email},
			"Je verpakkingen zijn binnen!",
			"collection-confirmation",
			map[string]interface{}{},
		)
	}

	j.SetLowerBound(now)
}
