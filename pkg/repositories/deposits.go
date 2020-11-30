package repositories

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/models"
)

type DepositsRepository struct {
	DB *gorm.DB
}

func (r DepositsRepository) ListCustomerDepositPayouts() ([]entities.CustomerDepositPayout, error) {
	ps := []models.CustomerDepositPayout{}
	err := r.DB.
		Preload("ReceivingUser").
		Preload("PayingUser").
		Preload("Return").
		Preload("Return.Package").
		Preload("Return.Package.PackageType").
		Find(ps).
		Error
	if err != nil {
		return nil, err
	}

	payouts := []entities.CustomerDepositPayout{}
	for _, p := range ps {
		payouts = append(payouts, p.ModelToEntity())
	}

	return payouts, nil
}

func (r DepositsRepository) ListUnpaidCustomerDepositPayouts() ([]entities.CustomerDepositPayout, error) {
	ps := []models.CustomerDepositPayout{}
	err := r.DB.
		Where("payment_id IS NULL OR payment_id = ''").
		Preload("ReceivingUser").
		Preload("PayingUser").
		Preload("Return").
		Preload("Return.Package").
		Preload("Return.Package.PackageType").
		Find(&ps).
		Error
	if err != nil {
		return nil, err
	}

	payouts := []entities.CustomerDepositPayout{}
	for _, p := range ps {
		payouts = append(payouts, p.ModelToEntity())
	}

	return payouts, nil
}

func (r DepositsRepository) ListCustomersDepositPayouts(userID uint) ([]entities.CustomerDepositPayout, error) {
	ps := []models.CustomerDepositPayout{}
	err := r.DB.
		Where("receiving_user_id = ?", userID).
		Preload("ReceivingUser").
		Preload("PayingUser").
		Preload("Return").
		Preload("Return.Package").
		Preload("Return.Package.PackageType").
		Find(&ps).
		Error
	if err != nil {
		return nil, err
	}

	payouts := []entities.CustomerDepositPayout{}
	for _, p := range ps {
		payouts = append(payouts, p.ModelToEntity())
	}

	return payouts, nil
}

func (r DepositsRepository) ListUnpaidCustomersDepositPayouts(userID uint) ([]entities.CustomerDepositPayout, error) {
	ps := []models.CustomerDepositPayout{}
	err := r.DB.
		Where("receiving_user_id = ? AND (payment_id IS NULL OR payment_id = '')", userID).
		Preload("ReceivingUser").
		Preload("PayingUser").
		Preload("Return").
		Preload("Return.Package").
		Preload("Return.Package.PackageType").
		Find(&ps).
		Error
	if err != nil {
		return nil, err
	}

	payouts := []entities.CustomerDepositPayout{}
	for _, p := range ps {
		payouts = append(payouts, p.ModelToEntity())
	}

	return payouts, nil
}

func (r DepositsRepository) CreateCustomerDepositPayout(receivingUserID, returnID uint) (*entities.CustomerDepositPayout, error) {
	p := models.CustomerDepositPayout{
		ReceivingUserID: receivingUserID,
		ReturnID:        returnID,
	}
	err := r.DB.
		Select("id", "created_at", "updated_at", "receiving_user_id", "return_id").
		Create(&p).
		Error
	if err != nil {
		return nil, err
	}

	err = r.DB.
		Model(&p).
		Preload("ReceivingUser").
		Preload("PayingUser").
		Preload("Return").
		Preload("Return.Package").
		Preload("Return.Package.PackageType").
		Find(&p).
		Error
	if err != nil {
		return nil, err
	}

	payout := p.ModelToEntity()

	return &payout, nil
}

func (r DepositsRepository) UpdateCustomerDepositPayout(payoutID, payingUserId uint, paymentID string) (*entities.CustomerDepositPayout, error) {
	p := models.CustomerDepositPayout{}
	err := r.DB.
		Model(&p).
		Where("id = ?", payoutID).
		Updates(&models.CustomerDepositPayout{
			PayingUserID: payingUserId,
			PaymentID:    paymentID,
		}).
		Error
	if err != nil {
		return nil, err
	}

	err = r.DB.
		Model(&p).
		Preload("ReceivingUser").
		Preload("PayingUser").
		Preload("Return").
		Preload("Return.Package").
		Preload("Return.Package.PackageType").
		Find(&p).
		Error
	if err != nil {
		return nil, err
	}

	payout := p.ModelToEntity()

	return &payout, nil
}
