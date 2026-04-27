package search

import (
	"unicode"

	"golang.org/x/text/language"
)

type Tokenizer struct {
	locales []language.Tag
}

func NewTokenizer() Tokenizer {
	return Tokenizer{locales}
}

func Tokenize(text string) []string {
}

func getScript(r rune) *unicode.RangeTable {
	//
}
