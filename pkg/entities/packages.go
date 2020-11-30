package entities

import uuid "github.com/satori/go.uuid"

type Package struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	PackageCode uuid.UUID   `json:"package_code"`
	PackageType PackageType `json:"package_type"`
}

type PackageType struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Manufacturer  string  `json:"manufacturer"`
	ProductCode   string  `json:"product_code"`
	DepositAmount float64 `json:"deposit_amount"`
}
