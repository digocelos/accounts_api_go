package account

import "time"

type Account struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Document  string    `json:"document"`
	Email     *string   `json:"email,omitempty"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateInput struct {
	Name     string  `json:"name"`
	Document string  `json:"document"`
	Email    *string `json:"email,omitempty"`
}

type UpdateInput struct {
	Name            *string `json:"name,omitempty"`
	Email           *string `json:"email,omitempty"`
	ExpectedVersion int     `json:"expected_version"`
}
