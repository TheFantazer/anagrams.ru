package dictionary

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadFromFile загружает словарь из текстового файла в Trie
// Формат файла: одно слово на строку
func LoadFromFile(filepath string) (*Trie, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open dictionary file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close dictionary file: %w", cerr)
		}
	}()

	trie := NewTrie()
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		word := strings.TrimSpace(scanner.Text())

		if word == "" || strings.HasPrefix(word, "#") {
			continue
		}

		trie.Insert(word)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading dictionary file: %w", err)
	}

	return trie, nil
}

// LoadMultipleDictionaries Загружает словари для нескольких языков
func LoadMultipleDictionaries(dictPaths map[string]string) (map[string]*Trie, error) {
	dictionaries := make(map[string]*Trie)

	for lang, path := range dictPaths {
		trie, err := LoadFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load %s dictionary: %w", lang, err)
		}
		dictionaries[lang] = trie
	}

	return dictionaries, nil
}
