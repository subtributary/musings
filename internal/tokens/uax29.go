package tokens

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/clipperhouse/uax29/v2/words"
)

type UAX29 struct{}

func (t UAX29) Tokens(text string) []string {
	tokens := words.FromString(text)
	tokens.Joiners(&words.Joiners[string]{
		Leading: []rune("#"),
		Middle:  []rune("\u200b"),
	})

	var results []string
	for tokens.Next() {
		token := tokens.Value()

		if !isWordToken(token) {
			continue
		}

		token = stripIgnoredChars(token)

		if len(token) == 0 {
			continue
		}

		results = append(results, token)
	}
	return results
}

func isWordToken(token string) bool {
	r, _ := utf8.DecodeRuneInString(token)
	return unicode.In(r, unicode.Letter, unicode.Number) || r == '#'
}

// stripIgnoredChars implements UAX #29's WB4 rule.
func stripIgnoredChars(token string) string {
	// Based on WB4, which is effectively M | FF9E | FF9F | (Cf except ZWSP).
	// We also ignore ZWSP because it adds no meaning.
	isIgnored := func(r rune) bool {
		return unicode.In(r, unicode.M, unicode.Cf) ||
			r == '\uff9e' || r == '\uff9f'
	}

	var builder strings.Builder
	for _, r := range token {
		if !isIgnored(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
