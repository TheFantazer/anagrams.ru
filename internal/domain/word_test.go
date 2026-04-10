package domain

import "testing"

func TestScoreWord(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		expected int
	}{
		{"empty word", "", 0},
		{"1 letter", "а", 0},
		{"2 letters", "он", 0},
		{"3 letters", "сон", 100},
		{"4 letters", "стол", 400},
		{"5 letters", "слово", 1200},
		{"6 letters", "словом", 2000},
		{"7 letters", "словами", 3000},
		{"8+ letters", "словарный", 0},
		{"uppercase", "СОН", 100},
		{"mixed case", "СлОвО", 1200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ScoreWord(tt.word)
			if got != tt.expected {
				t.Errorf("ScoreWord(%q) = %d, want %d", tt.word, got, tt.expected)
			}
		})
	}
}

func TestNewWord(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		expectedScore int
	}{
		{"3 letter word", "кот", 100},
		{"5 letter word", "слово", 1200},
		{"7 letter word", "словами", 3000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			word := NewWord(tt.text)
			if word.Text != tt.text {
				t.Errorf("NewWord(%q).Text = %q, want %q", tt.text, word.Text, tt.text)
			}
			if word.Score != tt.expectedScore {
				t.Errorf("NewWord(%q).Score = %d, want %d", tt.text, word.Score, tt.expectedScore)
			}
		})
	}
}
