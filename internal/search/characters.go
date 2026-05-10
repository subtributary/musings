package search

import (
	"slices"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// Character represents one or more runes that make up a symbol.
type Character string

// IsBreak returns true if the character causes a word break.
// It is the union of `!c.IswordPart()` and some end-of-word characters.
func (c Character) IsBreak() bool {
	return !c.IsWordPart()
}

// IsSkippable returns true if the character should be excluded from tokens.
func (c Character) IsSkippable() bool {
	r := c.primaryRune()
	if _, ok := slices.BinarySearch(forceSkips, r); ok {
		return true
	}
	if _, ok := slices.BinarySearch(forceNoSkips, r); ok {
		return false
	}
	return !(unicode.IsLetter(r) || unicode.IsNumber(r))
}

var forceSkips = []rune{
	'\uff9e', '\uff9f', // katakana sound marks
}

var forceNoSkips = []rune{
	'\u0027', // apostrophe
	'\u005f', // underscore
	'\u059e', // accent gershayim
}

// IsWordPart returns true if the character is part of a word or identifier.
func (c Character) IsWordPart() bool {
	r := c.primaryRune()
	if unicode.IsLetter(r) || unicode.IsNumber(r) {
		return true
	}
	if unicode.Is(unicode.Cf, r) && r != 0x070f && r != 0x8203 {
		return true
	}
	if unicode.Is(unicode.M, r) {
		return true
	}
	_, ok := slices.BinarySearch(extraLetters, r)
	return ok
}

var extraLetters = []rune{
	'\u0027', // apostrophe
	'\u005f', // underscore
	'\u00a0', // nbsp
	'\u00b7', // middle dot
	'\u00b8', // cedilla

	// modifier tone letters (Sk)
	'\u02c2', '\u02c3', '\u02c3', '\u02c4', '\u02c5', '\u02d2', '\u02d3',
	'\u02d4', '\u02d5', '\u02d6', '\u02d7', '\u02de', '\u02df', '\u02e5',
	'\u02e6', '\u02e7', '\u02e8', '\u02e9', '\u02ea', '\u02eb', '\u02ed',
	'\u02ef', '\u02f0', '\u02f1', '\u02f2', '\u02f3', '\u02f4', '\u02f5',
	'\u02f6', '\u02f7', '\u02f8', '\u02f9', '\u02fa', '\u02fb', '\u02fc',
	'\u02fd', '\u02fe', '\u02ff',

	'\u055a', '\u055b', '\u055c', '\u055e', '\u058a', // Armenian
	'\u059e', '\u05f3', // Hebrew
	'\u070f', // Syriac

	// modifier tone letters (Sk)
	'\ua708', '\ua709', '\ua70a', '\ua70b', '\ua70c', '\ua70d', '\ua70e', '\ua70f',
	'\ua710', '\ua711', '\ua712', '\ua713', '\ua714', '\ua715', '\ua716',
	'\ua720', '\ua721', '\ua789', '\ua78a', '\uab5b',

	'\uff9e', '\uff9f', // katakana sound marks
}

// Script returns the script that the character is part of.
func (c Character) Script() *unicode.RangeTable {
	r := c.primaryRune()

	if unicode.Is(unicode.Common, r) {
		return nil
	}

	for _, script := range unicode.Scripts {
		if unicode.Is(script, r) {
			return script
		}
	}

	return nil
}

func (c Character) primaryRune() rune {
	r, _ := utf8.DecodeRuneInString(string(c))
	return r
}

type CharReader struct {
	reader  norm.Iter
	current Character
}

func NewCharReader(text string) *CharReader {
	r := &CharReader{}
	r.reader.InitString(norm.NFC, text)
	return r
}

func (r *CharReader) Current() Character {
	return r.current
}

// Done returns true if there are no more characters after the current one.
func (r *CharReader) Done() bool {
	return r.reader.Done()
}

// Read reads the next character. Current can be used to get it.
func (r *CharReader) Read() Character {
	r.current = Character(r.reader.Next())
	return r.current
}
