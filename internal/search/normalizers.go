package search

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// NormalizeText is a language-agnostic normalizer.
// It should be applied to all text consumed by the search algorithms.
// The other `Normalize*` functions herein call this function.
//
// These transforms are applied:
//   - Unicode characters are canonicalized per NFKC.
//   - Whitespace is trimmed.
//   - Text is lowercased.
//
// The result is not always accurate—for example, the Turkish 'I' is
// incorrectly lowercased—but it is sufficient for equivalency checks.
// It should not be used for display purposes.
func NormalizeText(text string) string {
	text = norm.NFKC.String(text)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	return text
}

// NormalizeTag normalizes common ways of writing a tag into a canonical form.
// It should be applied to all tags consumed by the search algorithms.
// The result should not be used for display purposes.
func NormalizeTag(tag string) string {
	tag = NormalizeText(tag)
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, tag)
}
