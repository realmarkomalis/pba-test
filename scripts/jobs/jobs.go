package jobs

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"gitlab.com/markomalis/packback-api/pkg/services"
)

type Job struct {
	MailService services.IEmailService
	LowerBound  time.Time
	DB          *gorm.DB
}

func NewJob(DB *gorm.DB, f func(*Job)) *Job {
	es := services.NewEmailService(
		viper.GetString("email_service.domain"),
		viper.GetString("email_service.secret_key"),
	)

	j := Job{
		MailService: es,
		LowerBound:  time.Now(),
		DB:          DB,
	}

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minutes().Do(f, &j)
	s.StartAsync()

	return &j
}

func (j *Job) SetLowerBound(t time.Time) {
	j.LowerBound = t
}

func InitJobs(DB *gorm.DB) {
	NewJob(DB, SendDropOffCollectionNotifications)
	NewJob(DB, SendDropOffIntentNotifications)
	NewJob(DB, SendPickUpCollectionNotifications)
	NewJob(DB, SendPickupRequestNotifications)
	NewJob(DB, SendRefundConfirmationNotifications)
}
