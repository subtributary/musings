package search

import (
	"slices"
	"testing"
)

func TestTokenize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "lowercases Latin words",
			text: "Hello, WORLD!",
			want: []string{"hello", "world"},
		},
		{
			name: "keeps numbers",
			text: "Version 0.3 ships in 2026.",
			want: []string{"version", "0.3", "ships", "in", "2026"},
		},
		{
			name: "handles apostrophes",
			text: "Don't stop believing.",
			want: []string{"don't", "stop", "believing"},
		},
		{
			name: "normalizes compatibility forms",
			text: "Ｆｕｌｌｗｉｄｔｈ",
			want: []string{"fullwidth"},
		},
		{
			name: "keeps hebrew words",
			text: "שלום עולם",
			want: []string{"שלום", "עולם"},
		},
		{
			name: "keeps arabic words",
			text: "مرحبا بالعالم",
			want: []string{"مرحبا", "بالعالم"},
		},
		{
			name: "keeps hangul words as whole tokens",
			text: "한국어를 배웠다",
			want: []string{"한국어를", "배웠다"},
		},
		{
			name: "expands han into 1-grams",
			text: "北京大学",
			want: []string{"北", "京", "大", "学"},
		},
		{
			name: "expands kana into bigrams and trigrams",
			text: "カタカナ",
			want: []string{"カタ", "タカ", "カナ", "カタカ", "タカナ"},
		},
		{
			name: "returns short han token unchanged",
			text: "北",
			want: []string{"北"},
		},
		{
			name: "handles mixed scripts",
			text: "I visited 서울 and 北京.",
			want: []string{"i", "visited", "서울", "and", "北", "京"},
		},
		{
			name: "drops punctuation only",
			text: "!!! --- ...",
			want: nil,
		},
		{
			name: "empty string",
			text: "",
			want: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Tokenize(tt.text)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("Tokenize(%q): got %#v, want %#v", tt.text, got, tt.want)
			}
		})
	}
}
