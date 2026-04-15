package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewResult(t *testing.T) {
	validSessionID := uuid.New()
	foundWords := []string{"сон", "нос"}

	tests := []struct {
		name              string
		sessionID         uuid.UUID
		playerName        string
		playerFingerprint string
		foundWords        []string
		durationMs        int
		wantErr           error
		expectedName      string
	}{
		{
			name:              "valid result with name",
			sessionID:         validSessionID,
			playerName:        "Vasya",
			playerFingerprint: "fp123",
			foundWords:        foundWords,
			durationMs:        45000,
			wantErr:           nil,
			expectedName:      "Vasya",
		},
		{
			name:              "valid result without name - should default to Anonymous",
			sessionID:         validSessionID,
			playerName:        "",
			playerFingerprint: "fp123",
			foundWords:        foundWords,
			durationMs:        45000,
			wantErr:           nil,
			expectedName:      "Anonymous",
		},
		{
			name:              "invalid session ID",
			sessionID:         uuid.Nil,
			playerName:        "Vasya",
			playerFingerprint: "fp123",
			foundWords:        foundWords,
			durationMs:        45000,
			wantErr:           ErrInvalidSessionID,
		},
		{
			name:              "missing fingerprint",
			sessionID:         validSessionID,
			playerName:        "Vasya",
			playerFingerprint: "",
			foundWords:        foundWords,
			durationMs:        45000,
			wantErr:           ErrMissingFingerprint,
		},
		{
			name:              "invalid duration - zero",
			sessionID:         validSessionID,
			playerName:        "Vasya",
			playerFingerprint: "fp123",
			foundWords:        foundWords,
			durationMs:        0,
			wantErr:           ErrInvalidDuration,
		},
		{
			name:              "invalid duration - negative",
			sessionID:         validSessionID,
			playerName:        "Vasya",
			playerFingerprint: "fp123",
			foundWords:        foundWords,
			durationMs:        -100,
			wantErr:           ErrInvalidDuration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewResult(tt.sessionID, nil, tt.playerName, tt.playerFingerprint, tt.foundWords, tt.durationMs)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("NewResult() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("NewResult() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewResult() unexpected error = %v", err)
				return
			}

			if result.ID == uuid.Nil {
				t.Error("NewResult() ID should not be nil")
			}
			if result.SessionID != tt.sessionID {
				t.Errorf("NewResult() SessionID = %v, want %v", result.SessionID, tt.sessionID)
			}
			if result.PlayerName != tt.expectedName {
				t.Errorf("NewResult() PlayerName = %v, want %v", result.PlayerName, tt.expectedName)
			}
			if result.PlayerFingerprint != tt.playerFingerprint {
				t.Errorf("NewResult() PlayerFingerprint = %v, want %v", result.PlayerFingerprint, tt.playerFingerprint)
			}
			if len(result.FoundWords) != len(tt.foundWords) {
				t.Errorf("NewResult() FoundWords length = %v, want %v", len(result.FoundWords), len(tt.foundWords))
			}
			if result.WordCount != len(tt.foundWords) {
				t.Errorf("NewResult() WordCount = %v, want %v", result.WordCount, len(tt.foundWords))
			}
			if result.DurationMs != tt.durationMs {
				t.Errorf("NewResult() DurationMs = %v, want %v", result.DurationMs, tt.durationMs)
			}

			expectedScore := CalculateScore(tt.foundWords)
			if result.Score != expectedScore {
				t.Errorf("NewResult() Score = %v, want %v", result.Score, expectedScore)
			}
		})
	}
}

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		want  int
	}{
		{
			name:  "empty list",
			words: []string{},
			want:  0,
		},
		{
			name:  "single word",
			words: []string{"сон"},
			want:  100,
		},
		{
			name:  "multiple words same length",
			words: []string{"сон", "нос", "сор"},
			want:  300, // 100 * 3
		},
		{
			name:  "multiple words different lengths",
			words: []string{"сон", "слово", "словами"},
			want:  100 + 1200 + 2800, // 3+5+7 letters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateScore(tt.words)
			if got != tt.want {
				t.Errorf("CalculateScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_ValidateWords(t *testing.T) {
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
		name       string
		foundWords []string
		wantErr    error
	}{
		{
			name:       "all valid words",
			foundWords: []string{"сон", "нос"},
			wantErr:    nil,
		},
		{
			name:       "contains invalid word",
			foundWords: []string{"сон", "дом"},
			wantErr:    ErrInvalidWord,
		},
		{
			name:       "empty list",
			foundWords: []string{},
			wantErr:    nil,
		},
		{
			name:       "only invalid words",
			foundWords: []string{"дом", "кот"},
			wantErr:    ErrInvalidWord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &Result{
				ID:                uuid.New(),
				SessionID:         session.ID,
				PlayerName:        "Test",
				PlayerFingerprint: "fp123",
				FoundWords:        tt.foundWords,
				WordCount:         len(tt.foundWords),
				Score:             CalculateScore(tt.foundWords),
				DurationMs:        45000,
				PlayedAt:          time.Now().UTC(),
			}

			err := result.ValidateWords(session)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidateWords() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err != tt.wantErr {
					t.Errorf("ValidateWords() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateWords() unexpected error = %v", err)
			}
		})
	}
}

func TestResult_CalculateAccuracy(t *testing.T) {
	tests := []struct {
		name             string
		sessionWords     []string
		foundWords       []string
		expectedAccuracy float64
	}{
		{
			name:             "found all words",
			sessionWords:     []string{"сон", "нос", "сор"},
			foundWords:       []string{"сон", "нос", "сор"},
			expectedAccuracy: 100.0,
		},
		{
			name:             "found half words",
			sessionWords:     []string{"сон", "нос", "сор", "рот"},
			foundWords:       []string{"сон", "нос"},
			expectedAccuracy: 50.0,
		},
		{
			name:             "found no words",
			sessionWords:     []string{"сон", "нос", "сор"},
			foundWords:       []string{},
			expectedAccuracy: 0.0,
		},
		{
			name:             "found one of three",
			sessionWords:     []string{"сон", "нос", "сор"},
			foundWords:       []string{"сон"},
			expectedAccuracy: 100.0 / 3.0,
		},
		{
			name:             "session has no valid words",
			sessionWords:     []string{},
			foundWords:       []string{},
			expectedAccuracy: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &Session{
				ID:          uuid.New(),
				Letters:     "сонрто",
				Language:    "ru",
				TimeLimit:   60,
				LetterCount: 6,
				ValidWords:  tt.sessionWords,
				MaxScore:    calculateMaxScore(tt.sessionWords),
				CreatedAt:   time.Now().UTC(),
			}

			result := &Result{
				ID:                uuid.New(),
				SessionID:         session.ID,
				PlayerName:        "Test",
				PlayerFingerprint: "fp123",
				FoundWords:        tt.foundWords,
				WordCount:         len(tt.foundWords),
				Score:             CalculateScore(tt.foundWords),
				DurationMs:        45000,
				PlayedAt:          time.Now().UTC(),
			}

			got := result.CalculateAccuracy(session)

			// Используем epsilon для сравнения float
			epsilon := 0.0000001
			if (got-tt.expectedAccuracy) > epsilon || (tt.expectedAccuracy-got) > epsilon {
				t.Errorf("CalculateAccuracy() = %v, want %v", got, tt.expectedAccuracy)
			}
		})
	}
}
