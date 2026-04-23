package domain

import (
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID  `json:"id"`
	Letters     string     `json:"letters"`
	Language    string     `json:"language"`
	TimeLimit   int        `json:"time_limit"`
	LetterCount int        `json:"letter_count"`
	ValidWords  []string   `json:"valid_words"`
	MaxScore    int        `json:"max_score"`
	CreatorID   *uuid.UUID `json:"creator_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	HideLetters bool       `json:"hide_letters"`
}

func NewSession(letters, language string, timeLimit, letterCount int, validWords []string, hideLetters bool) (*Session, error) {
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
		ID:          uuid.New(),
		Letters:     letters,
		Language:    language,
		TimeLimit:   timeLimit,
		LetterCount: letterCount,
		ValidWords:  validWords,
		MaxScore:    maxScore,
		CreatedAt:   time.Now().UTC(),
		HideLetters: hideLetters,
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
