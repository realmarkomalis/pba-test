package entities

type DepositItem struct {
	Status     string  `json:"status"`
	StatusCode int     `json:"status_code"`
	Package    Package `json:"package"`
	Return     Return  `json:"return"`
}

type CustomerDepositPayout struct {
	ID            uint   `json:"id"`
	ReceivingUser User   `json:"receiving_user"`
	PayingUser    User   `json:"paying_user"`
	Return        Return `json:"return"`
	PaymentID     string `json:"payment_id"`
}
