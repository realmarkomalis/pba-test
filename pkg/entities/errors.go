package entities

import "fmt"

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Link    string `json:"link"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("API Error: %s, for more info visit: %s", e.Message, e.Link)
}
