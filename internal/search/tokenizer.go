package search

import (
	"iter"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

type Tokenizer interface {
	Tokens() iter.Seq[string]
}

type DocumentTokenizer struct {
	reader *CharReader
}

func NewDocumentTokenizer(text string) DocumentTokenizer {
	return DocumentTokenizer{
		reader: NewCharReader(text),
	}
}

func (t *DocumentTokenizer) Tokens() iter.Seq[string] {
	return func(yield func(string) bool) {
		for !t.reader.Done() {
			consumeUntil(t.reader, func(ch Character) bool {
				return ch.IsWordPart()
			})

			for value := range t.subtokens(consumeWord(t.reader)) {
				value = norm.NFKC.String(value)
				value = strings.ToLower(value)
				if value != "" {
					if !yield(value) {
						return
					}
				}
			}
		}
	}
}

func (t *DocumentTokenizer) subtokens(text string) iter.Seq[string] {
	r, _ := utf8.DecodeRuneInString(text)
	switch Character(r).Script() {
	case unicode.Han:
		tok := NewNGramTokenizer(text, 1, 1)
		return tok.Tokens()
	case unicode.Hiragana, unicode.Katakana:
		tok := NewNGramTokenizer(text, 2, 3)
		return tok.Tokens()
	default:
		return slices.Values([]string{text})
	}
}

// NGramTokenizer tokenizes text into n-grams
type NGramTokenizer struct {
	minN, maxN int
	runes      []rune
}

func NewNGramTokenizer(text string, minN, maxN int) NGramTokenizer {
	return NGramTokenizer{minN, maxN, []rune(text)}
}

func (t *NGramTokenizer) Tokens() iter.Seq[string] {
	// Text smaller than the minimum is just returned.
	if len(t.runes) < t.minN {
		result := []string{string(t.runes)}
		return slices.Values(result)
	}

	// maxN can't be longer than the text.
	maxN := min(t.maxN, len(t.runes))

	// Return overlapping chunks of sizes from minN to maxN.
	return func(yield func(string) bool) {
		for n := t.minN; n <= maxN; n++ {
			for i := 0; i+n <= len(t.runes); i++ {
				if !yield(string(t.runes[i : i+n])) {
					return
				}
			}
		}
	}
}

// QueryTokenizer is the same as DocumentTokenizer except it supports tags.
type QueryTokenizer struct {
	DocumentTokenizer
}

func NewQueryTokenizer(text string) QueryTokenizer {
	return QueryTokenizer{NewDocumentTokenizer(text)}
}

func (t *QueryTokenizer) Tokens() iter.Seq[string] {
	// This should be nearly identical to DocumentTokenizer.Tokens.
	// The only difference should be the tag check.
	return func(yield func(string) bool) {
		for !t.reader.Done() {
			consumeUntil(t.reader, func(ch Character) bool {
				return ch.IsWordPart() || ch == "#"
			})

			if value := consumeTag(t.reader); value != "" {
				value = norm.NFKC.String(value)
				value = strings.ToLower(value)
				if !yield(value) {
					return
				}
				continue
			}

			for value := range t.subtokens(consumeWord(t.reader)) {
				value = norm.NFKC.String(value)
				value = strings.ToLower(value)
				if value != "" {
					if !yield(value) {
						return
					}
				}
			}
		}
	}
}

func consumeTag(reader *CharReader) (result string) {
	if reader.Current() != "#" {
		return
	}

	builder := strings.Builder{}
	for ch := reader.Read(); ch.IsWordPart(); ch = reader.Read() {
		if !ch.IsSkippable() {
			builder.WriteString(string(ch))
		}
		if reader.Done() {
			break
		}
	}

	if builder.Len() > 0 {
		result = "#" + builder.String()
	}

	return
}

// consumeUntil consumes characters until a condition is met or end of stream.
// The reader is positioned at the first character meeting the condition.
func consumeUntil(reader *CharReader, until func(ch Character) bool) {
	for ; !reader.Done(); reader.Read() {
		if until(reader.Current()) {
			break
		}
	}
}

// consumeWord consumes word characters starting at the current character and
// returns their string ending at the first non-word character.
func consumeWord(reader *CharReader) string {
	builder := strings.Builder{}
	var prevScript *unicode.RangeTable = nil
	for ch := reader.Current(); ch.IsWordPart(); ch = reader.Read() {
		// Break on script change
		script := ch.Script()
		if script != nil {
			if prevScript != nil && prevScript != script {
				break
			}
			prevScript = script
		}

		if !ch.IsSkippable() {
			builder.WriteString(string(ch))
		}
		if reader.Done() {
			break
		}
	}
	return builder.String()
}
