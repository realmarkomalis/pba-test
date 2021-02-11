package entities

type Supplier struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Users        []User        `json:"users"`
	PackageTypes []PackageType `json:"package_types"`
}
