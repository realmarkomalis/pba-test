package entities

import uuid "github.com/satori/go.uuid"

type Package struct {
	ID          uint      `json:"id"`
	PackageCode uuid.UUID `json:"package_code"`
	Name        string    `json:"name"`
}
