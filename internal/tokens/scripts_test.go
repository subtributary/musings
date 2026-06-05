package tokens_test

import (
	"slices"
	"testing"

	"github.com/subtributary/musings/internal/tokens"
)

func TestScripts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want []tokens.ScriptToken
	}{
		{
			name: "empty",
			text: "",
			want: nil,
		},
		{
			name: "single",
			text: "abc",
			want: []tokens.ScriptToken{
				{"Latin", "abc"},
			},
		},
		{
			name: "multi",
			text: "a술晴",
			want: []tokens.ScriptToken{
				{"Latin", "a"},
				{"Hangul", "술"},
				{"Han", "晴"},
			},
		},
		{
			name: "modifiers",
			text: "c\u00b8",
			want: []tokens.ScriptToken{
				{"Latin", "c\u00b8"},
			},
		},
		{
			name: "common",
			text: "! ?",
			want: nil,
		},
		{
			name: "leading common",
			text: "!abc",
			want: []tokens.ScriptToken{
				{"Latin", "abc"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			subject := tokens.Scripts{}
			got := subject.Tokens(tt.text)

			if !slices.Equal(tt.want, got) {
				t.Errorf("Tokens returned %v, want %v", got, tt.want)
			}
		})
	}
}
