package jobs

import (
	"fmt"
	"time"

	"gitlab.com/markomalis/packback-api/pkg/models"
)

func SendDropOffCollectionNotifications(j *Job) {
	fmt.Println("SendDropOffCollectionNotifications")
	now := time.Now().UTC()
	users := []models.User{}

	err := j.DB.
		Select("users.email").
		Joins("JOIN drop_off_intents on users.id = drop_off_intents.user_id").
		Joins("JOIN drop_off_collects on drop_off_intents.return_id = drop_off_collects.return_id").
		Where("drop_off_collects.created_at BETWEEN ? AND ?", j.LowerBound, now).
		Group("users.email").
		Find(&users).
		Error
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, u := range users {
		fmt.Printf("SendDropOffCollectionNotifications mail user with email: %s/n", u.Email)
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
