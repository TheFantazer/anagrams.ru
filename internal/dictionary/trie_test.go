package dictionary

import (
	"reflect"
	"sort"
	"testing"
)

func TestNewTrie(t *testing.T) {
	trie := NewTrie()
	if trie == nil {
		t.Fatal("NewTrie() returned nil")
	}
	if trie.root == nil {
		t.Fatal("NewTrie() root is nil")
	}
	if trie.root.children == nil {
		t.Fatal("NewTrie() root.children is nil")
	}
}

func TestTrie_Insert(t *testing.T) {
	trie := NewTrie()

	words := []string{"cat", "car", "card", "care", "careful"}
	for _, word := range words {
		trie.Insert(word)
	}

	// Проверяем что все слова можно найти
	for _, word := range words {
		if !trie.Search(word) {
			t.Errorf("Trie.Search(%q) = false after Insert, want true", word)
		}
	}
}

func TestTrie_Insert_CaseInsensitive(t *testing.T) {
	trie := NewTrie()

	trie.Insert("CAT")

	tests := []string{"cat", "Cat", "CAT", "cAt"}
	for _, word := range tests {
		if !trie.Search(word) {
			t.Errorf("Trie.Search(%q) = false, want true (case insensitive)", word)
		}
	}
}

func TestTrie_Search(t *testing.T) {
	trie := NewTrie()
	words := []string{"cat", "car", "card"}

	for _, word := range words {
		trie.Insert(word)
	}

	tests := []struct {
		name string
		word string
		want bool
	}{
		{"existing word 1", "cat", true},
		{"existing word 2", "car", true},
		{"existing word 3", "card", true},
		{"non-existing word", "dog", false},
		{"prefix of existing word", "ca", false},
		{"extension of existing word", "cards", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trie.Search(tt.word)
			if got != tt.want {
				t.Errorf("Trie.Search(%q) = %v, want %v", tt.word, got, tt.want)
			}
		})
	}
}

func TestTrie_FindAllWords(t *testing.T) {
	trie := NewTrie()

	// Русские слова
	words := []string{"сон", "нос", "сор", "рот", "сто", "он"}
	for _, word := range words {
		trie.Insert(word)
	}

	tests := []struct {
		name    string
		letters string
		want    []string
	}{
		{
			name:    "найти несколько слов",
			letters: "сонрто",
			want:    []string{"сон", "нос", "сор", "рот", "сто"},
		},
		{
			name:    "найти несколько слов из трёх букв",
			letters: "сон",
			want:    []string{"сон", "нос"},
		},
		{
			name:    "нет подходящих слов",
			letters: "абв",
			want:    []string{},
		},
		{
			name:    "пустые буквы",
			letters: "",
			want:    []string{},
		},
		{
			name:    "буквы в разном регистре",
			letters: "СоНрТо",
			want:    []string{"сон", "нос", "сор", "рот", "сто"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trie.FindAllWords(tt.letters)

			// Сортируем для корректного сравнения
			sort.Strings(got)
			sort.Strings(tt.want)

			// Нормализуем пустые слайсы для сравнения
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trie.FindAllWords(%q) = %v, want %v", tt.letters, got, tt.want)
			}
		})
	}
}

func TestTrie_FindAllWords_MinLength(t *testing.T) {
	trie := NewTrie()

	// Добавляем слова разной длины
	words := []string{"а", "он", "сон", "слон"}
	for _, word := range words {
		trie.Insert(word)
	}

	// FindAllWords должен возвращать только слова >= 3 букв
	got := trie.FindAllWords("аслон")

	// Ожидаем только "сон" и "слон", слова "а" и "он" должны быть отфильтрованы
	want := []string{"сон", "слон"}

	sort.Strings(got)
	sort.Strings(want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Trie.FindAllWords() returned words shorter than 3 letters: got %v, want %v", got, want)
	}
}

func TestTrie_FindAllWords_DuplicateLetters(t *testing.T) {
	trie := NewTrie()

	words := []string{"мама", "папа"}
	for _, word := range words {
		trie.Insert(word)
	}

	// "мама" требует 2x 'м' и 2x 'а'
	tests := []struct {
		name    string
		letters string
		want    []string
	}{
		{
			name:    "достаточно букв для 'мама'",
			letters: "ммааппп",
			want:    []string{"мама", "папа"},
		},
		{
			name:    "недостаточно 'а' для 'мама'",
			letters: "ммп",
			want:    []string{},
		},
		{
			name:    "только одна 'м' - не хватит для 'мама'",
			letters: "мапа",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trie.FindAllWords(tt.letters)

			sort.Strings(got)
			sort.Strings(tt.want)

			// Нормализуем пустые слайсы для сравнения
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trie.FindAllWords(%q) = %v, want %v", tt.letters, got, tt.want)
			}
		})
	}
}

func TestTrie_Size(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		want  int
	}{
		{
			name:  "empty trie",
			words: []string{},
			want:  0,
		},
		{
			name:  "single word",
			words: []string{"cat"},
			want:  1,
		},
		{
			name:  "multiple words",
			words: []string{"cat", "car", "card", "care"},
			want:  4,
		},
		{
			name:  "duplicate words - should count as one",
			words: []string{"cat", "cat", "dog"},
			want:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := NewTrie()
			for _, word := range tt.words {
				trie.Insert(word)
			}

			got := trie.Size()
			if got != tt.want {
				t.Errorf("Trie.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
