package usecases

import (
	"gitlab.com/markomalis/packback-api/pkg/entities"
	"gitlab.com/markomalis/packback-api/pkg/repositories"
)

type DepositsUsecase struct {
	ReturnsRepository        repositories.IReturnRepository
	DropOffIntentsRepository repositories.DropOffIntentsRepository
	DepositsRepository       repositories.DepositsRepository
}

func (u DepositsUsecase) ListUnpaidCustomerDepositPayouts() ([]entities.CustomerDepositPayout, error) {
	ps, err := u.DepositsRepository.ListUnpaidCustomerDepositPayouts()
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (u DepositsUsecase) GetUserDepositItems(userID uint) ([]entities.DepositItem, error) {
	di, err := u.getUserDepositItems(userID)
	if err != nil {
		return nil, err
	}

	ps, err := u.DepositsRepository.ListCustomersDepositPayouts(userID)
	if err != nil {
		return nil, err
	}

	items := []entities.DepositItem{}
	items = append(items, u.openReturns(ps, di)...)
	items = append(items, u.claimedReturns(ps, di)...)
	items = append(items, u.paidReturns(ps, di)...)

	return items, nil
}

func (u DepositsUsecase) ClaimCustomerDeposits(userID uint) ([]entities.CustomerDepositPayout, error) {
	payouts, err := u.DepositsRepository.ListCustomersDepositPayouts(userID)
	if err != nil {
		return nil, err
	}

	deposits, err := u.getUserDepositItems(userID)
	if err != nil {
		return nil, err
	}

	collectedReturns := []entities.Return{}
	for _, d := range deposits {
		if d.Return.Status == "Collected" || d.Return.Status == "Fulfilled" {
			collectedReturns = append(collectedReturns, d.Return)
		}
	}

	claimableReturns := u.claimableReturns(payouts, collectedReturns)
	ps := []entities.CustomerDepositPayout{}
	for _, r := range claimableReturns {
		p, err := u.DepositsRepository.CreateCustomerDepositPayout(userID, r.ID)
		if err == nil {
			ps = append(ps, *p)
		}
	}

	return ps, nil
}

func (u DepositsUsecase) PayCustomerDeposits(receivingUserID, payingUserID uint, paymentID string) (ps []entities.CustomerDepositPayout, e error) {
	payouts, err := u.DepositsRepository.ListUnpaidCustomersDepositPayouts(receivingUserID)
	if err != nil {
		return nil, err
	}

	for _, p := range payouts {
		cp, err := u.DepositsRepository.UpdateCustomerDepositPayout(p.ID, payingUserID, paymentID)
		if err == nil {
			ps = append(ps, *cp)
		}
	}

	return ps, nil
}

func (u DepositsUsecase) getUserDepositItems(userID uint) ([]entities.DepositItem, error) {
	items := []entities.DepositItem{}
	rs, err := u.ReturnsRepository.GetReturnRequests(userID)
	if err != nil {
		return nil, err
	}

	for _, r := range rs {
		r, err := u.ReturnsRepository.GetReturn(r.ID)
		if err == nil {
			items = append(items, entities.DepositItem{
				Status:     r.Status,
				StatusCode: r.StatusCode,
				Package:    r.Package,
				Return:     *r,
			})
		}
	}

	ds, err := u.DropOffIntentsRepository.GetUserDropOffIntents(userID)
	if err != nil {
		return nil, err
	}

	for _, d := range ds {
		r, err := u.ReturnsRepository.GetReturn(d.ReturnID)
		if err == nil {
			items = append(items, entities.DepositItem{
				Status:     r.Status,
				StatusCode: r.StatusCode,
				Package:    r.Package,
				Return:     *r,
			})
		}
	}

	return items, nil
}

func (u DepositsUsecase) claimableReturns(ps []entities.CustomerDepositPayout, rs []entities.Return) (cr []entities.Return) {
	m := make(map[uint]bool)

	for _, p := range ps {
		m[p.Return.ID] = true
	}

	for _, r := range rs {
		if _, ok := m[r.ID]; !ok {
			cr = append(cr, r)
		}
	}
	return
}

func (u DepositsUsecase) paidReturns(ps []entities.CustomerDepositPayout, rs []entities.DepositItem) (cr []entities.DepositItem) {
	m := make(map[uint]bool)

	for _, p := range ps {
		if p.PaymentID != "" {
			m[p.Return.ID] = true
		}
	}

	for _, r := range rs {
		if _, ok := m[r.Return.ID]; ok {
			r.Status = "Paid"
			cr = append(cr, r)
		}
	}
	return
}

func (u DepositsUsecase) claimedReturns(ps []entities.CustomerDepositPayout, rs []entities.DepositItem) (cr []entities.DepositItem) {
	m := make(map[uint]bool)

	for _, p := range ps {
		if p.PaymentID == "" {
			m[p.Return.ID] = true
		}
	}

	for _, r := range rs {
		if _, ok := m[r.Return.ID]; ok {
			r.Status = "Claimed"
			cr = append(cr, r)
		}
	}
	return
}

func (u DepositsUsecase) openReturns(ps []entities.CustomerDepositPayout, rs []entities.DepositItem) (cr []entities.DepositItem) {
	m := make(map[uint]bool)

	for _, p := range ps {
		m[p.Return.ID] = true
	}

	for _, r := range rs {
		if _, ok := m[r.Return.ID]; !ok {
			cr = append(cr, r)
		}
	}
	return
}
