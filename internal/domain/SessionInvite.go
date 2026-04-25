package domain

import (
	"time"

	"github.com/google/uuid"
)

type SessionInvite struct {
	ID         uuid.UUID `json:"id"`
	SessionID  uuid.UUID `json:"session_id"`
	FromUserID uuid.UUID `json:"from_user_id"`
	ToUserID   uuid.UUID `json:"to_user_id"`
	Status     string    `json:"status"` // "pending", "accepted", "declined"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewSessionInvite(sessionID, fromUserID, toUserID uuid.UUID) *SessionInvite {
	now := time.Now().UTC()
	return &SessionInvite{
		ID:         uuid.New(),
		SessionID:  sessionID,
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
