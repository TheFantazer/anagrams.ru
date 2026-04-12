package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Email         *string   `json:"email,omitempty"`
	Password      string    `json:"-"` // Never expose password in JSON
	OAuthProvider *string   `json:"oauth_provider,omitempty"`
	OAuthID       *string   `json:"oauth_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewUser(username, email, password string) (*User, error) {
	if len(username) < 3 || len(username) > 30 {
		return nil, ErrInvalidUsername
	}

	if password != "" && len(password) < 6 {
		return nil, ErrPasswordTooShort
	}

	var emailPtr *string
	if email != "" {
		emailPtr = &email
	}

	return &User{
		ID:        uuid.New(),
		Username:  username,
		Email:     emailPtr,
		Password:  password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
