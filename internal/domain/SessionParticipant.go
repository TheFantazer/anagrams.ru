package domain

import (
	"time"

	"github.com/google/uuid"
)

// SessionParticipant represents a participant in a game session
type SessionParticipant struct {
	ID        uuid.UUID  `json:"id"`
	SessionID uuid.UUID  `json:"session_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Role      string     `json:"role"` // "creator" or "opponent"
	JoinedAt  time.Time  `json:"joined_at"`
	StartedAt *time.Time `json:"started_at,omitempty"` // NULL if hasn't started the game yet
}

// NewSessionParticipant creates a new session participant
func NewSessionParticipant(sessionID, userID uuid.UUID, role string) *SessionParticipant {
	now := time.Now().UTC()
	return &SessionParticipant{
		ID:        uuid.New(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      role,
		JoinedAt:  now,
		StartedAt: nil,
	}
}

// MarkAsStarted marks the participant as having started the game
func (p *SessionParticipant) MarkAsStarted() {
	now := time.Now().UTC()
	p.StartedAt = &now
}

// HasStarted returns whether the participant has started the game
func (p *SessionParticipant) HasStarted() bool {
	return p.StartedAt != nil
}
