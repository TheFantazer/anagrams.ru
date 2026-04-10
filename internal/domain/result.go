package domain

import (
	"time"

	"github.com/google/uuid"
)

type Result struct {
	ID                uuid.UUID `json:"id"`
	SessionID         uuid.UUID `json:"session_id"`
	PlayerName        string    `json:"player_name"`
	PlayerFingerprint string    `json:"player_fingerprint"`
	FoundWords        []string  `json:"found_words"`
	WordCount         int       `json:"word_count"`
	Score             int       `json:"score"`
	DurationMs        int       `json:"duration_ms"`
	PlayedAt          time.Time `json:"played_at"`
}

func NewResult(sessionID uuid.UUID, playerName, playerFingerprint string, foundWords []string, durationMs int) (*Result, error) {
	if sessionID == uuid.Nil {
		return nil, ErrInvalidSessionID
	}
	if playerFingerprint == "" {
		return nil, ErrMissingFingerprint
	}
	if durationMs <= 0 {
		return nil, ErrInvalidDuration
	}

	if playerName == "" {
		playerName = "Anonymous"
	}

	score := CalculateScore(foundWords)

	return &Result{
		ID:                uuid.New(),
		SessionID:         sessionID,
		PlayerName:        playerName,
		PlayerFingerprint: playerFingerprint,
		FoundWords:        foundWords,
		WordCount:         len(foundWords),
		Score:             score,
		DurationMs:        durationMs,
		PlayedAt:          time.Now().UTC(),
	}, nil
}

// CalculateScore вычисляет общий счёт из списка слов
func CalculateScore(words []string) int {
	totalScore := 0
	for _, word := range words {
		totalScore += ScoreWord(word)
	}
	return totalScore
}

// ValidateWords проверяет слова против сессии
func (r *Result) ValidateWords(session *Session) error {
	for _, word := range r.FoundWords {
		if !session.IsValid(word) {
			return ErrInvalidWord
		}
	}
	return nil
}

// CalculateAccuracy вычисляет процент найденных слов
func (r *Result) CalculateAccuracy(session *Session) float64 {
	if len(session.ValidWords) == 0 {
		return 0
	}
	return float64(r.WordCount) / float64(len(session.ValidWords)) * 100
}
