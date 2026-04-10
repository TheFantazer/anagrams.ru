package dictionary

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	// Создаём временную директорию для тестов
	tmpDir := t.TempDir()

	t.Run("valid dictionary file", func(t *testing.T) {
		// Создаём временный файл словаря
		dictFile := filepath.Join(tmpDir, "test_dict.txt")
		content := `cat
dog
bird
# комментарий
fish

elephant`
		err := os.WriteFile(dictFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		trie, err := LoadFromFile(dictFile)
		if err != nil {
			t.Fatalf("LoadFromFile() error = %v, want nil", err)
		}

		// Проверяем что все слова загружены (кроме комментария и пустой строки)
		expectedWords := []string{"cat", "dog", "bird", "fish", "elephant"}
		for _, word := range expectedWords {
			if !trie.Search(word) {
				t.Errorf("Trie.Search(%q) = false after LoadFromFile, want true", word)
			}
		}

		// Проверяем размер
		if trie.Size() != len(expectedWords) {
			t.Errorf("Trie.Size() = %v, want %v", trie.Size(), len(expectedWords))
		}
	})

	t.Run("file with only comments and empty lines", func(t *testing.T) {
		dictFile := filepath.Join(tmpDir, "empty_dict.txt")
		content := `# comment 1
# comment 2


# another comment`
		err := os.WriteFile(dictFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		trie, err := LoadFromFile(dictFile)
		if err != nil {
			t.Fatalf("LoadFromFile() error = %v, want nil", err)
		}

		if trie.Size() != 0 {
			t.Errorf("Trie.Size() = %v, want 0 for empty dictionary", trie.Size())
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := LoadFromFile("non_existent_file.txt")
		if err == nil {
			t.Error("LoadFromFile() with non-existent file should return error")
		}
	})

	t.Run("russian words", func(t *testing.T) {
		dictFile := filepath.Join(tmpDir, "ru_dict.txt")
		content := `сон
нос
дом
кот
мир`
		err := os.WriteFile(dictFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		trie, err := LoadFromFile(dictFile)
		if err != nil {
			t.Fatalf("LoadFromFile() error = %v, want nil", err)
		}

		expectedWords := []string{"сон", "нос", "дом", "кот", "мир"}
		for _, word := range expectedWords {
			if !trie.Search(word) {
				t.Errorf("Trie.Search(%q) = false, want true", word)
			}
		}
	})

	t.Run("words with extra whitespace", func(t *testing.T) {
		dictFile := filepath.Join(tmpDir, "whitespace_dict.txt")
		content := `  cat
   dog
bird
  fish  `
		err := os.WriteFile(dictFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		trie, err := LoadFromFile(dictFile)
		if err != nil {
			t.Fatalf("LoadFromFile() error = %v, want nil", err)
		}

		// Whitespace должен быть обрезан
		expectedWords := []string{"cat", "dog", "bird", "fish"}
		for _, word := range expectedWords {
			if !trie.Search(word) {
				t.Errorf("Trie.Search(%q) = false, want true (whitespace should be trimmed)", word)
			}
		}
	})
}

func TestLoadMultipleDictionaries(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("load multiple valid dictionaries", func(t *testing.T) {
		// Создаём английский словарь
		enFile := filepath.Join(tmpDir, "en.txt")
		enContent := `cat
dog
bird`
		err := os.WriteFile(enFile, []byte(enContent), 0644)
		if err != nil {
			t.Fatalf("failed to create en file: %v", err)
		}

		// Создаём русский словарь
		ruFile := filepath.Join(tmpDir, "ru.txt")
		ruContent := `кот
собака
птица`
		err = os.WriteFile(ruFile, []byte(ruContent), 0644)
		if err != nil {
			t.Fatalf("failed to create ru file: %v", err)
		}

		dictPaths := map[string]string{
			"en": enFile,
			"ru": ruFile,
		}

		dictionaries, err := LoadMultipleDictionaries(dictPaths)
		if err != nil {
			t.Fatalf("LoadMultipleDictionaries() error = %v, want nil", err)
		}

		// Проверяем что загружено 2 словаря
		if len(dictionaries) != 2 {
			t.Errorf("LoadMultipleDictionaries() loaded %v dictionaries, want 2", len(dictionaries))
		}

		// Проверяем английский словарь
		enTrie, ok := dictionaries["en"]
		if !ok {
			t.Fatal("English dictionary not found")
		}
		if !enTrie.Search("cat") {
			t.Error("English dictionary doesn't contain 'cat'")
		}

		// Проверяем русский словарь
		ruTrie, ok := dictionaries["ru"]
		if !ok {
			t.Fatal("Russian dictionary not found")
		}
		if !ruTrie.Search("кот") {
			t.Error("Russian dictionary doesn't contain 'кот'")
		}
	})

	t.Run("fail on missing dictionary file", func(t *testing.T) {
		dictPaths := map[string]string{
			"en": filepath.Join(tmpDir, "en.txt"),
			"ru": "non_existent_ru.txt",
		}

		_, err := LoadMultipleDictionaries(dictPaths)
		if err == nil {
			t.Error("LoadMultipleDictionaries() should return error when file is missing")
		}
	})

	t.Run("empty dictionary map", func(t *testing.T) {
		dictPaths := map[string]string{}

		dictionaries, err := LoadMultipleDictionaries(dictPaths)
		if err != nil {
			t.Fatalf("LoadMultipleDictionaries() error = %v, want nil", err)
		}

		if len(dictionaries) != 0 {
			t.Errorf("LoadMultipleDictionaries() with empty map should return empty map, got %v entries", len(dictionaries))
		}
	})
}
