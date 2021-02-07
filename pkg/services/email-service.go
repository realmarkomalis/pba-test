package services

import (
	"context"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type IEmailService interface {
	SendEmail(from string, to []string, subject, body string) (string, error)
	SendTemplateEmail(from string, to []string, subject, template string, variables map[string]interface{}) (string, error)
}

type EmailService struct {
	domain    string
	secretKey string
}

func NewEmailService(domain, secretKey string) IEmailService {
	return EmailService{
		domain:    domain,
		secretKey: secretKey,
	}
}

func (es EmailService) SendEmail(from string, to []string, subject, body string) (string, error) {
	mg := mailgun.NewMailgun(es.domain, es.secretKey)
	mg.SetAPIBase(mailgun.APIBaseEU)
	m := mg.NewMessage(
		from, subject, body, strings.Join(to, ","),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (es EmailService) SendTemplateEmail(from string, to []string, subject, template string, variables map[string]interface{}) (string, error) {
	mg := mailgun.NewMailgun(es.domain, es.secretKey)
	mg.SetAPIBase(mailgun.APIBaseEU)
	m := mg.NewMessage(
		from, subject, "", strings.Join(to, ","),
	)
	m.SetTemplate(template)

	for key, element := range variables {
		m.AddTemplateVariable(key, element)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	if err != nil {
		return "", err
	}

	return id, nil
}
