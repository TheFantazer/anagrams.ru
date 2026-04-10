package dictionary

import "strings"

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}
}

func (t *Trie) Insert(word string) {
	word = strings.ToLower(word)
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

func (t *Trie) Search(word string) bool {
	word = strings.ToLower(word)
	node := t.root

	for _, ch := range word {
		if node.children[ch] == nil {
			return false
		}
		node = node.children[ch]
	}
	return node.isEnd
}

func (t *Trie) FindAllWords(letters string) []string {
	letters = strings.ToLower(letters)

	letterCount := make(map[rune]int)
	for _, ch := range letters {
		letterCount[ch]++
	}

	var result []string
	var current []rune

	t.dfs(t.root, letterCount, current, &result)

	return result
}

func (t *Trie) dfs(node *TrieNode, letterCount map[rune]int, current []rune, result *[]string) {
	// Если нашли конец слова, добавляем его
	if node.isEnd && len(current) >= 3 { // минимум 3 буквы
		*result = append(*result, string(current))
	}

	for ch, count := range letterCount {
		if count > 0 && node.children[ch] != nil {
			letterCount[ch]--
			current = append(current, ch)

			t.dfs(node.children[ch], letterCount, current, result)

			letterCount[ch]++
			current = current[:len(current)-1]
		}
	}
}

func (t *Trie) Size() int {
	return t.countWords(t.root)
}

func (t *Trie) countWords(node *TrieNode) int {
	count := 0
	if node.isEnd {
		count++
	}
	for _, child := range node.children {
		count += t.countWords(child)
	}
	return count
}
