package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID `db:"id" json:"id"`
	Username           string    `db:"username" json:"username"`
	Email              *string   `db:"email" json:"email,omitempty"`
	Password           string    `db:"password" json:"-"`
	OAuthProvider      *string   `db:"oauth_provider" json:"oauth_provider,omitempty"`
	OAuthID            *string   `db:"oauth_id" json:"oauth_id,omitempty"`
	DefaultLetterCount int       `db:"default_letter_count" json:"default_letter_count"`
	DefaultLanguage    string    `db:"default_language" json:"default_language"`
	DefaultTimeLimit   int       `db:"default_time_limit" json:"default_time_limit"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
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
