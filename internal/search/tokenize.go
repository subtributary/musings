package search

import (
	"strings"
	"unicode"

	"github.com/clipperhouse/uax29/words"
	"golang.org/x/text/unicode/norm"
)

// Tokenize converts text into normalized search tokens.
//
// It uses Unicode word segmentation as the base tokenizer, then applies
// search-oriented normalization and script-specific expansion.
//
// Current behavior:
//   - UAX #29 word segmentation
//   - Unicode NFKC normalization
//   - lowercase normalization
//   - punctuation/space-only tokens are discarded
//   - Han/Kana tokens are expanded into character n-grams
//   - other scripts are emitted as whole tokens
//
// This intentionally does not perform stemming, lemmatization, language
// detection, Korean particle stripping, or Arabic/Hebrew morphology yet.
func Tokenize(text string) []string {
	if text == "" {
		return nil
	}

	segmenter := words.NewSegmenter([]byte(text))
	tokens := make([]string, 0)

	for segmenter.Next() {
		raw := string(segmenter.Bytes())
		token := normalizeToken(raw)
		if token == "" || !containsTokenRune(token) {
			continue
		}

		tokens = append(tokens, expandToken(token)...)
	}

	return tokens
}

// normalizeToken applies normalization that should be consistent between
// indexing and querying.
func normalizeToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}

	// NFKC folds many compatibility forms into their canonical-ish equivalents.
	// This is useful for search because visually similar / compatibility forms
	// should usually match.
	token = norm.NFKC.String(token)

	// Locale-neutral lowercasing. This is intentionally not perfect for every
	// language, but index and query normalization remain consistent.
	token = strings.ToLower(token)

	return token
}

// containsTokenRune reports whether token has at least one rune worth indexing.
// This filters punctuation-only and symbol-only segments emitted by UAX #29.
func containsTokenRune(token string) bool {
	for _, r := range token {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return true
		}
	}
	return false
}

func expandToken(token string) []string {
	script := dominantSearchScript(token)

	switch script {
	case scriptHan:
		// The tokenizer we use splits this into 1-grams already,
		// so this doesn't actually have any effect.
		// Maybe I'll fix this later.
		return ngrams(token, 1, 2)
	case scriptKana:
		return ngrams(token, 2, 3)
	default:
		return []string{token}
	}
}

type searchScript int

const (
	scriptOther searchScript = iota
	scriptLatin
	scriptGreek
	scriptCyrillic
	scriptHebrew
	scriptArabic
	scriptHangul
	scriptHan
	scriptKana
)

// dominantSearchScript returns the first recognized script with a letter rune.
// UAX #29 has already produced word-ish tokens, so this is mainly used to decide
// whether a token needs n-gram expansion.
func dominantSearchScript(token string) searchScript {
	for _, r := range token {
		switch {
		case unicode.In(r, unicode.Han):
			return scriptHan
		case unicode.In(r, unicode.Hiragana, unicode.Katakana):
			return scriptKana
		case unicode.In(r, unicode.Hangul):
			return scriptHangul
		case unicode.In(r, unicode.Arabic):
			return scriptArabic
		case unicode.In(r, unicode.Hebrew):
			return scriptHebrew
		case unicode.In(r, unicode.Latin):
			return scriptLatin
		case unicode.In(r, unicode.Greek):
			return scriptGreek
		case unicode.In(r, unicode.Cyrillic):
			return scriptCyrillic
		}
	}

	return scriptOther
}

// ngrams returns character n-grams for a token.
func ngrams(token string, min, max int) []string {
	runes := []rune(token)

	if len(runes) < min {
		return []string{token}
	}

	if max > len(runes) {
		max = len(runes)
	}

	// Preallocate roughly (not exact, but close enough)
	capEstimate := 0
	for n := min; n <= max; n++ {
		capEstimate += len(runes) - n + 1
	}
	out := make([]string, 0, capEstimate)

	for n := min; n <= max; n++ {
		for i := 0; i+n <= len(runes); i++ {
			out = append(out, string(runes[i:i+n]))
		}
	}

	return out
}
