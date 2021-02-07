package jobs

import (
	"fmt"
	"time"

	"gitlab.com/markomalis/packback-api/pkg/models"
)

func SendRefundConfirmationNotifications(j *Job) {
	fmt.Println("SendRefundConfirmationNotifications")
	now := time.Now().UTC()
	users := []models.User{}

	err := j.DB.
		Select("users.email").
		Joins("JOIN customer_deposit_payouts on users.id = customer_deposit_payouts.receiving_user_id").
		Where("customer_deposit_payouts.created_at BETWEEN ? AND ?", j.LowerBound, now).
		Group("users.email").
		Find(&users).
		Error
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, u := range users {
		fmt.Printf("SendRefundConfirmationNotifications mail user with email: %s/n", u.Email)
		j.MailService.SendTemplateEmail(
			"\"PackBack\" <info@packback.app>",
			[]string{u.Email},
			"Uitbetaling van je statiegeld",
			"refund-deposit",
			map[string]interface{}{},
		)
	}

	j.SetLowerBound(now)
}
