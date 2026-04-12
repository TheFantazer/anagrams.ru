package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	Email              *string   `json:"email,omitempty"`
	Password           string    `json:"-"`
	OAuthProvider      *string   `json:"oauth_provider,omitempty"`
	OAuthID            *string   `json:"oauth_id,omitempty"`
	DefaultLetterCount int       `json:"default_letter_count"`
	DefaultLanguage    string    `json:"default_language"`
	DefaultTimeLimit   int       `json:"default_time_limit"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
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
		ID:                 uuid.New(),
		Username:           username,
		Email:              emailPtr,
		Password:           password,
		DefaultLetterCount: 7,
		DefaultLanguage:    "ru",
		DefaultTimeLimit:   60,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}, nil
}
