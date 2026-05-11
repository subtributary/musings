package search_test

import (
	"testing"
	"unicode"

	"github.com/subtributary/musings/internal/search"
)

func TestCharacterIsBreak(t *testing.T) {
	t.Parallel()

	// We mostly care about characters that are `IsWordPart() == true`
	// because the inverse is the first check that `IsBreak` does.
	tests := []struct {
		name  string
		chars string
		want  bool
	}{
		{
			name:  "not word part",
			chars: `~!@$%^&*()+-=/.,?><\|]}[{`,
			want:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for _, c := range tt.chars {
				subject := search.Character(c)
				got := subject.IsBreak()
				if got != tt.want {
					t.Fatalf("%#x.IsBreak(): got = %v, want %v", c, got, tt.want)
				}
			}
		})
	}
}

func TestCharacterIsSkippable(t *testing.T) {
	t.Parallel()

	// We only care about characters that are `IsWordPart() == true`
	// because no others will be checked for skipping by our tokenizer.
	tests := []struct {
		name  string
		chars string
		want  bool
	}{
		{
			name:  "basic letters and numbers",
			chars: "09Az",
			want:  false,
		},
		{
			name:  "spaces that do not break",
			chars: "\u00a0",
			want:  true,
		},
		{
			name:  "ambiguous symbols that are skipped",
			chars: "\u00b7",
			want:  true,
		},
		{
			name:  "ambiguous symbols that are not skipped",
			chars: "\u0027\u059e",
			want:  false,
		},
		{
			name:  "joiners",
			chars: "\u200d",
			want:  true,
		},
		{
			// I want to skip all modifiers left after NFC normalization.
			name:  "Mc | Me | Mn | FF9E | FF9F | Cf",
			chars: "\u302f\u0327\u0903\uff9e\uff9f\u00ad",
			want:  true,
		},
		{
			name: "uax29 extra letters",
			chars: "\u00b8\u02c2\u02c3\u02c3\u02c4\u02c5" +
				"\u02d2\u02d3\u02d4\u02d5\u02d6\u02d7" +
				"\u02de\u02df\u02e5\u02e6\u02e7\u02e8\u02e9\u02ea\u02eb\u02ed" +
				"\u02ef\u02f0\u02f1\u02f2\u02f3\u02f4\u02f5\u02f6\u02f7" +
				"\u02f8\u02f9\u02fa\u02fb\u02fc\u02fd\u02fe\u02ff" +
				"\u055a\u055b\u055c\u055e\u058a\u05f3\u070f" +
				"\ua708\ua709\ua70a\ua70b\ua70c\ua70d\ua70e\ua70f" +
				"\ua710\ua711\ua712\ua713\ua714\ua715\ua716" +
				"\ua720\ua721\ua789\ua78a\uab5b",
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for _, c := range tt.chars {
				subject := search.Character(c)
				got := subject.IsSkippable()
				if got != tt.want {
					t.Fatalf("%#x.IsSkippable(): got = %v, want %v", c, got, tt.want)
				}
			}
		})
	}
}

func TestCharacterScript(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		chars string
		want  *unicode.RangeTable
	}{
		{
			name:  "Han",
			chars: "北京大学",
			want:  unicode.Han,
		},
		{
			name:  "Hangul",
			chars: "안녕하세요",
			want:  unicode.Hangul,
		},
		{
			name:  "Hebrew",
			chars: "עִבְרִי",
			want:  unicode.Hebrew,
		},
		{
			name:  "Hiragana",
			chars: "こんにちは",
			want:  unicode.Hiragana,
		},
		{
			name:  "Katakana",
			chars: "コンニチハ",
			want:  unicode.Katakana,
		},
		{
			name:  "Latin",
			chars: "AbcDef",
			want:  unicode.Latin,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for _, c := range tt.chars {
				subject := search.Character(c)
				got := subject.Script()
				if got != tt.want {
					t.Errorf("%#v.Script(): wrong", c)
				}
			}
		})
	}
}

func TestCharacterIsWordPart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		chars string
		want  bool
	}{
		{
			name:  "basic_letters_and_numbers",
			chars: "09Az",
			want:  true,
		},
		{
			name:  "basic punctuation",
			chars: "!?",
			want:  false,
		},
		{
			name: "spaces that break",
			chars: "" +
				"\r\n\t\f\u0085\u2028\u2029" + // Newlines
				"\u0020\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007" +
				"\u2008\u2009\u200a\u202f\u205f\u3000",
			want: false,
		},
		{
			name:  "spaces that do not break",
			chars: "\u00a0",
			want:  true,
		},
		{
			name:  "joiners",
			chars: "\u200d",
			want:  true,
		},
		{
			name: "ambiguous symbols that break",
			chars: "" +
				"\u2018\u2019\uff07" + // similar to apostrophe (U+0027)
				"\u002e\u2024\ufe52\uff0e" + // full stop and similar
				"\u003a\ufe13\ufe55\uff1a" + // Swedish colon + similar
				"\u0387\u2027" + // similar to middle dot (U+00B7)
				"\u05f4\u0022", // punctuation gershayim + similar
			want: false,
		},
		{
			name: "ambiguous symbols that do not break",
			chars: "" +
				"\u0027" + // apostrophe
				"\u00b7" + // middle dot
				"\u059e", // accent gershayim
			want: true,
		},
		{
			// This is a simplification of UAX #29's WB4 rule for our use case.
			name:  "Mc | Me | Mn | FF9E | FF9F | (Cf except ZWSP)",
			chars: "\u302f\u0327\u0903\uff9e\uff9f\u00ad",
			want:  true,
		},
		{
			name: "uax29 extra letters",
			chars: "\u00b8\u02c2\u02c3\u02c3\u02c4\u02c5" +
				"\u02d2\u02d3\u02d4\u02d5\u02d6\u02d7" +
				"\u02de\u02df\u02e5\u02e6\u02e7\u02e8\u02e9\u02ea\u02eb\u02ed" +
				"\u02ef\u02f0\u02f1\u02f2\u02f3\u02f4\u02f5\u02f6\u02f7" +
				"\u02f8\u02f9\u02fa\u02fb\u02fc\u02fd\u02fe\u02ff" +
				"\u055a\u055b\u055c\u055e\u058a\u05f3\u070f" +
				"\ua708\ua709\ua70a\ua70b\ua70c\ua70d\ua70e\ua70f" +
				"\ua710\ua711\ua712\ua713\ua714\ua715\ua716" +
				"\ua720\ua721\ua789\ua78a\uab5b",
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for _, c := range tt.chars {
				subject := search.Character(c)
				got := subject.IsWordPart()
				if got != tt.want {
					t.Fatalf("%#x.IsWordPart(): got = %v, want %v", c, got, tt.want)
				}
			}
		})
	}
}
