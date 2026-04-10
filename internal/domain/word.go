package domain

import "unicode/utf8"

type Word struct {
	Text  string `json:"text"`
	Score int    `json:"score"`
}

func ScoreWord(word string) int {
	length := utf8.RuneCountInString(word)

	switch {
	case length <= 2:
		return 0
	case length == 3:
		return 100
	case length == 4:
		return 400
	case length == 5:
		return 1200
	case length == 6:
		return 2000
	case length == 7:
		return 3000
	default:
		return 0
	}
}

func NewWord(text string) Word {
	return Word{
		Text:  text,
		Score: ScoreWord(text),
	}
}
