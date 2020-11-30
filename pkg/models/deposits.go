package models

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
)

type CustomerDepositPayout struct {
	gorm.Model
	ReceivingUser   User
	ReceivingUserID uint
	PayingUser      User
	PayingUserID    uint
	Return          Return
	ReturnID        uint
	PaymentID       string
}

func (p CustomerDepositPayout) ModelToEntity() entities.CustomerDepositPayout {
	return entities.CustomerDepositPayout{
		ID:            p.ID,
		ReceivingUser: p.ReceivingUser.ModelToEntity(),
		PayingUser:    p.PayingUser.ModelToEntity(),
		Return:        p.Return.ModelToEntity(),
		PaymentID:     p.PaymentID,
	}
}
