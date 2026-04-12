package handler

import (
	"time"

	"github.com/google/uuid"
)

type CreateSessionRequest struct {
	Language    string `json:"language"`
	LetterCount int    `json:"letter_count"`
	TimeLimit   int    `json:"time_limit"`
}

type SessionResponse struct {
	ID          uuid.UUID `json:"id"`
	Letters     string    `json:"letters"`
	Language    string    `json:"language"`
	TimeLimit   int       `json:"time_limit"`
	LetterCount int       `json:"letter_count"`
	MaxScore    int       `json:"max_score"`
	ValidWords  []string  `json:"valid_words"`
	CreatedAt   time.Time `json:"created_at"`
}

type SubmitResultRequest struct {
	PlayerName  string   `json:"player_name"`
	Fingerprint string   `json:"fingerprint"`
	FoundWords  []string `json:"found_words"`
	DurationMs  int      `json:"duration_ms"`
}

type ResultResponse struct {
	ID         uuid.UUID `json:"id"`
	SessionID  uuid.UUID `json:"session_id"`
	PlayerName string    `json:"player_name"`
	WordCount  int       `json:"word_count"`
	Score      int       `json:"score"`
	DurationMs int       `json:"duration_ms"`
	PlayedAt   time.Time `json:"played_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
