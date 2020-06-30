package entities

type Restaurant struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	User User   `json:"user"`
}
