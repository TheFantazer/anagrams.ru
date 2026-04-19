package domain

import (
	"time"

	"github.com/google/uuid"
)

type SessionInvite struct {
	ID        uuid.UUID `json:"id"`
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewSessionInvite(sessionID, userID uuid.UUID) *SessionInvite {
	return &SessionInvite{
		ID:        uuid.New(),
		SessionID: sessionID,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}
}
