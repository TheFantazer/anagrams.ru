package dictionary

import (
	"math/rand"
	"strings"
	"time"
)

type LetterGenerator struct {
	rng *rand.Rand
}

func NewLetterGenerator() *LetterGenerator {
	return &LetterGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

var russianLetterFrequency = map[rune]int{
	'а': 80, 'б': 16, 'в': 45, 'г': 17, 'д': 30, 'е': 85, 'ё': 1,
	'ж': 10, 'з': 17, 'и': 74, 'й': 12, 'к': 35, 'л': 44, 'м': 32,
	'н': 67, 'о': 110, 'п': 28, 'р': 47, 'с': 55, 'т': 63, 'у': 26,
	'ф': 3, 'х': 10, 'ц': 5, 'ч': 15, 'ш': 7, 'щ': 4, 'ъ': 1,
	'ы': 19, 'ь': 18, 'э': 3, 'ю': 6, 'я': 20,
}

var englishLetterFrequency = map[rune]int{
	'a': 82, 'b': 15, 'c': 28, 'd': 43, 'e': 127, 'f': 22, 'g': 20,
	'h': 61, 'i': 70, 'j': 2, 'k': 8, 'l': 40, 'm': 24, 'n': 67,
	'o': 75, 'p': 19, 'q': 1, 'r': 60, 's': 63, 't': 91, 'u': 28,
	'v': 10, 'w': 24, 'x': 2, 'y': 20, 'z': 1,
}

func (lg *LetterGenerator) GenerateLetters(language string, count int) string {
	var frequency map[rune]int

	switch language {
	case "ru":
		frequency = russianLetterFrequency
	case "en":
		frequency = englishLetterFrequency
	default:
		frequency = russianLetterFrequency
	}

	// Создаём weighted pool букв
	var letterPool []rune
	for letter, weight := range frequency {
		for i := 0; i < weight; i++ {
			letterPool = append(letterPool, letter)
		}
	}

	// Генерируем случайные буквы
	var result []rune
	for i := 0; i < count; i++ {
		randomIndex := lg.rng.Intn(len(letterPool))
		result = append(result, letterPool[randomIndex])
	}

	return string(result)
}

// GenerateBalancedLetters генерирует набор с гарантией наличия гласных
func (lg *LetterGenerator) GenerateBalancedLetters(language string, count int) string {
	vowels := map[string][]rune{
		"ru": {'а', 'е', 'и', 'о', 'у', 'ы', 'э', 'ю', 'я'},
		"en": {'a', 'e', 'i', 'o', 'u'},
	}

	letters := lg.GenerateLetters(language, count)
	runes := []rune(strings.ToLower(letters))

	hasVowel := false
	for _, r := range runes {
		for _, v := range vowels[language] {
			if r == v {
				hasVowel = true
				break
			}
		}
		if hasVowel {
			break
		}
	}

	if !hasVowel && len(runes) > 0 {
		vowelList := vowels[language]
		runes[lg.rng.Intn(len(runes))] = vowelList[lg.rng.Intn(len(vowelList))]
	}

	return string(runes)
}
