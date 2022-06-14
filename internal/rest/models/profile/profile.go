package profile

import (
	"time"
)

type Profile struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Mobile     string    `json:"mobile"`
	IsActive   bool      `json:"is_active"`
	Name       string    `json:"name"`
	MiddleName string    `json:"middle_name"`
	SureName   string    `json:"sure_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Roles      []string  `json:"roles"`
}
