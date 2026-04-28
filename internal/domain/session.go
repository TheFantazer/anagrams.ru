package domain

import (
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Session struct {
	ID            uuid.UUID  `json:"id"`
	Letters       string     `json:"letters"`
	Language      string     `json:"language"`
	TimeLimit     int        `json:"time_limit"`
	LetterCount   int        `json:"letter_count"`
	ValidWords    []string   `json:"valid_words"`
	MaxScore      int        `json:"max_score"`
	CreatorID     *uuid.UUID `json:"creator_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	MaxOpponents  int        `json:"max_opponents"` // How many opponents allowed (1 for current 1v1, N for future NvN)
	InviteMode    string     `json:"invite_mode"`   // "link" or "friend"
	IsDaily       bool       `json:"is_daily"`      // Whether this is a daily puzzle
	DailyPuzzleID *uuid.UUID `json:"daily_puzzle_id,omitempty"`
}

func NewSession(letters, language string, timeLimit, letterCount int, validWords []string) (*Session, error) {
	if letters == "" {
		return nil, ErrInvalidLetters
	}
	if utf8.RuneCountInString(letters) != letterCount {
		return nil, ErrLetterCountMismatch
	}
	if language != "en" && language != "ru" {
		return nil, ErrUnsupportedLanguage
	}
	if timeLimit <= 0 || timeLimit > 600 {
		return nil, ErrInvalidTimeLimit
	}
	if len(validWords) == 0 {
		return nil, ErrNoValidWords
	}
	maxScore := calculateMaxScore(validWords)
	return &Session{
		ID:           uuid.New(),
		Letters:      letters,
		Language:     language,
		TimeLimit:    timeLimit,
		LetterCount:  letterCount,
		ValidWords:   validWords,
		MaxScore:     maxScore,
		CreatedAt:    time.Now().UTC(),
		MaxOpponents: 1,      // Default to 1v1
		InviteMode:   "link", // Default to link mode
	}, nil
}

func (s *Session) IsExpired() bool {
	// Multiplayer sessions are valid for 7 days
	expiryDuration := 7 * 24 * time.Hour
	expiresAt := s.CreatedAt.Add(expiryDuration)
	return time.Now().UTC().After(expiresAt)
}
func (s *Session) IsValid(word string) bool {
	for _, ValidWord := range s.ValidWords {
		if ValidWord == word {
			return true
		}
	}
	return false
}

func calculateMaxScore(words []string) int {
	totalScore := 0
	for _, word := range words {
		totalScore += ScoreWord(word)
	}
	return totalScore
}
