package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewSession(t *testing.T) {
	validWords := []string{"сон", "нос", "сор"}

	tests := []struct {
		name        string
		letters     string
		language    string
		timeLimit   int
		letterCount int
		validWords  []string
		wantErr     error
	}{
		{
			name:        "valid session ru",
			letters:     "сонрто",
			language:    "ru",
			timeLimit:   60,
			letterCount: 6,
			validWords:  validWords,
			wantErr:     nil,
		},
		{
			name:        "valid session en",
			letters:     "abcdefg",
			language:    "en",
			timeLimit:   120,
			letterCount: 7,
			validWords:  []string{"cab", "bed"},
			wantErr:     nil,
		},
		{
			name:        "empty letters",
			letters:     "",
			language:    "ru",
			timeLimit:   60,
			letterCount: 0,
			validWords:  validWords,
			wantErr:     ErrInvalidLetters,
		},
		{
			name:        "letter count mismatch",
			letters:     "сон",
			language:    "ru",
			timeLimit:   60,
			letterCount: 5,
			validWords:  validWords,
			wantErr:     ErrLetterCountMismatch,
		},
		{
			name:        "unsupported language",
			letters:     "abcdefg",
			language:    "fr",
			timeLimit:   60,
			letterCount: 7,
			validWords:  validWords,
			wantErr:     ErrUnsupportedLanguage,
		},
		{
			name:        "invalid time limit - zero",
			letters:     "сонрто",
			language:    "ru",
			timeLimit:   0,
			letterCount: 6,
			validWords:  validWords,
			wantErr:     ErrInvalidTimeLimit,
		},
		{
			name:        "invalid time limit - too high",
			letters:     "сонрто",
			language:    "ru",
			timeLimit:   601,
			letterCount: 6,
			validWords:  validWords,
			wantErr:     ErrInvalidTimeLimit,
		},
		{
			name:        "no valid words",
			letters:     "сонрто",
			language:    "ru",
			timeLimit:   60,
			letterCount: 6,
			validWords:  []string{},
			wantErr:     ErrNoValidWords,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := NewSession(tt.letters, tt.language, tt.timeLimit, tt.letterCount, tt.validWords)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("NewSession() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("NewSession() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewSession() unexpected error = %v", err)
				return
			}

			if session.ID == uuid.Nil {
				t.Error("NewSession() ID should not be nil")
			}
			if session.Letters != tt.letters {
				t.Errorf("NewSession() Letters = %v, want %v", session.Letters, tt.letters)
			}
			if session.Language != tt.language {
				t.Errorf("NewSession() Language = %v, want %v", session.Language, tt.language)
			}
			if session.TimeLimit != tt.timeLimit {
				t.Errorf("NewSession() TimeLimit = %v, want %v", session.TimeLimit, tt.timeLimit)
			}
			if session.LetterCount != tt.letterCount {
				t.Errorf("NewSession() LetterCount = %v, want %v", session.LetterCount, tt.letterCount)
			}
			if len(session.ValidWords) != len(tt.validWords) {
				t.Errorf("NewSession() ValidWords length = %v, want %v", len(session.ValidWords), len(tt.validWords))
			}
			if session.MaxScore != calculateMaxScore(tt.validWords) {
				t.Errorf("NewSession() MaxScore = %v, want %v", session.MaxScore, calculateMaxScore(tt.validWords))
			}
		})
	}
}

func TestSession_IsExpired(t *testing.T) {
	validWords := []string{"test"}

	tests := []struct {
		name      string
		timeLimit int
		createdAt time.Time
		want      bool
	}{
		{
			name:      "not expired - just created",
			timeLimit: 60,
			createdAt: time.Now().UTC(),
			want:      false,
		},
		{
			name:      "expired - 7 minutes ago with 1 minute limit",
			timeLimit: 60,
			createdAt: time.Now().UTC().Add(-7 * time.Minute),
			want:      true,
		},
		{
			name:      "not expired - 30 seconds ago with 60 second limit",
			timeLimit: 60,
			createdAt: time.Now().UTC().Add(-30 * time.Second),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &Session{
				ID:          uuid.New(),
				Letters:     "test",
				Language:    "en",
				TimeLimit:   tt.timeLimit,
				LetterCount: 4,
				ValidWords:  validWords,
				MaxScore:    100,
				CreatedAt:   tt.createdAt,
			}

			got := session.IsExpired()
			if got != tt.want {
				t.Errorf("Session.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_IsValid(t *testing.T) {
	session := &Session{
		ID:          uuid.New(),
		Letters:     "сонрто",
		Language:    "ru",
		TimeLimit:   60,
		LetterCount: 6,
		ValidWords:  []string{"сон", "нос", "сор", "рот"},
		MaxScore:    400,
		CreatedAt:   time.Now().UTC(),
	}

	tests := []struct {
		name string
		word string
		want bool
	}{
		{"valid word 1", "сон", true},
		{"valid word 2", "нос", true},
		{"invalid word", "дом", false},
		{"empty word", "", false},
		{"partial match", "со", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := session.IsValid(tt.word)
			if got != tt.want {
				t.Errorf("Session.IsValid(%q) = %v, want %v", tt.word, got, tt.want)
			}
		})
	}
}

func TestCalculateMaxScore(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		want  int
	}{
		{
			name:  "single 3-letter word",
			words: []string{"сон"},
			want:  100,
		},
		{
			name:  "multiple words",
			words: []string{"сон", "нос", "сорт"},
			want:  100 + 100 + 400, // 3+3+4 letters
		},
		{
			name:  "empty list",
			words: []string{},
			want:  0,
		},
		{
			name:  "mixed lengths",
			words: []string{"сон", "слово", "словами"},
			want:  100 + 1200 + 2800, // 3+5+7 letters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateMaxScore(tt.words)
			if got != tt.want {
				t.Errorf("calculateMaxScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
