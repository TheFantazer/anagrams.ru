package dictionary

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestNewLetterGenerator(t *testing.T) {
	gen := NewLetterGenerator()
	if gen == nil {
		t.Fatal("NewLetterGenerator() returned nil")
	}
	if gen.rng == nil {
		t.Fatal("NewLetterGenerator() rng is nil")
	}
}

func TestLetterGenerator_GenerateLetters(t *testing.T) {
	gen := NewLetterGenerator()

	tests := []struct {
		name     string
		language string
		count    int
	}{
		{"russian 7 letters", "ru", 7},
		{"english 7 letters", "en", 7},
		{"russian 5 letters", "ru", 5},
		{"english 10 letters", "en", 10},
		{"default language (unknown)", "unknown", 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			letters := gen.GenerateLetters(tt.language, tt.count)

			// Проверяем что сгенерировано правильное количество букв
			actualCount := utf8.RuneCountInString(letters)
			if actualCount != tt.count {
				t.Errorf("GenerateLetters() generated %v letters, want %v", actualCount, tt.count)
			}

			// Проверяем что все символы в нижнем регистре
			if letters != strings.ToLower(letters) {
				t.Errorf("GenerateLetters() should return lowercase letters, got %q", letters)
			}
		})
	}
}

func TestLetterGenerator_GenerateLetters_Language(t *testing.T) {
	gen := NewLetterGenerator()

	t.Run("russian generates cyrillic", func(t *testing.T) {
		letters := gen.GenerateLetters("ru", 20)

		// Проверяем что все буквы кириллические
		for _, r := range letters {
			if !isCyrillic(r) {
				t.Errorf("GenerateLetters(ru) generated non-cyrillic character: %c", r)
			}
		}
	})

	t.Run("english generates latin", func(t *testing.T) {
		letters := gen.GenerateLetters("en", 20)

		// Проверяем что все буквы латинские
		for _, r := range letters {
			if !isLatin(r) {
				t.Errorf("GenerateLetters(en) generated non-latin character: %c", r)
			}
		}
	})
}

func TestLetterGenerator_GenerateLetters_Randomness(t *testing.T) {
	gen := NewLetterGenerator()

	// Генерируем несколько раз и проверяем что результаты разные
	results := make(map[string]bool)
	iterations := 10

	for i := 0; i < iterations; i++ {
		letters := gen.GenerateLetters("en", 7)
		results[letters] = true
	}

	// С высокой вероятностью должно быть несколько разных результатов
	if len(results) < 2 {
		t.Errorf("GenerateLetters() should produce different results, got %v unique values in %v iterations", len(results), iterations)
	}
}

func TestLetterGenerator_GenerateBalancedLetters(t *testing.T) {
	gen := NewLetterGenerator()

	tests := []struct {
		name     string
		language string
		count    int
	}{
		{"russian 7 letters", "ru", 7},
		{"english 7 letters", "en", 7},
		{"russian 5 letters", "ru", 5},
		{"english 10 letters", "en", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			letters := gen.GenerateBalancedLetters(tt.language, tt.count)

			// Проверяем что сгенерировано правильное количество букв
			actualCount := utf8.RuneCountInString(letters)
			if actualCount != tt.count {
				t.Errorf("GenerateBalancedLetters() generated %v letters, want %v", actualCount, tt.count)
			}

			// Проверяем что есть хотя бы одна гласная
			hasVowel := false
			vowels := getVowels(tt.language)

			for _, r := range letters {
				for _, v := range vowels {
					if r == v {
						hasVowel = true
						break
					}
				}
				if hasVowel {
					break
				}
			}

			if !hasVowel {
				t.Errorf("GenerateBalancedLetters(%q) = %q has no vowels", tt.language, letters)
			}
		})
	}
}

func TestLetterGenerator_GenerateBalancedLetters_AlwaysHasVowel(t *testing.T) {
	gen := NewLetterGenerator()

	// Запускаем много раз чтобы убедиться что всегда есть гласная
	for i := 0; i < 50; i++ {
		letters := gen.GenerateBalancedLetters("ru", 7)

		hasVowel := false
		vowels := getVowels("ru")

		for _, r := range letters {
			for _, v := range vowels {
				if r == v {
					hasVowel = true
					break
				}
			}
			if hasVowel {
				break
			}
		}

		if !hasVowel {
			t.Fatalf("GenerateBalancedLetters() iteration %v: %q has no vowels", i, letters)
		}
	}
}

func TestLetterGenerator_GenerateBalancedLetters_ZeroCount(t *testing.T) {
	gen := NewLetterGenerator()

	letters := gen.GenerateBalancedLetters("ru", 0)

	if len(letters) != 0 {
		t.Errorf("GenerateBalancedLetters() with count=0 should return empty string, got %q", letters)
	}
}

// Helper functions

func isCyrillic(r rune) bool {
	return (r >= 'а' && r <= 'я') || r == 'ё'
}

func isLatin(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func getVowels(language string) []rune {
	vowels := map[string][]rune{
		"ru": {'а', 'е', 'и', 'о', 'у', 'ы', 'э', 'ю', 'я'},
		"en": {'a', 'e', 'i', 'o', 'u'},
	}
	return vowels[language]
}
