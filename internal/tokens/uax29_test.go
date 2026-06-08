package tokens_test

import (
	"slices"
	"testing"

	"github.com/subtributary/musings/internal/tokens"
)

func TestUAX29(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "empty",
			text: "",
			want: nil,
		},
		{
			name: "Tibetan syllable delimiter",
			text: "དཔེ་ཆ",
			want: []string{"དཔ", "ཆ"},
		},
		{
			name: "hashtag",
			text: "#exid #weare ##",
			want: []string{"#exid", "#weare", "#", "#"},
		},

		/* from UAX #29 rules */
		{
			// Effectively M | FF9E | FF9F | (Cf except ZWSP)
			name: "WB4: Ignore extend and format characters",
			text: "\ufeff\u302f술 ZW\u200bSP ZW\u200cNJ ZW\u200dJ S\u00adHY",
			want: []string{"술", "ZWSP", "ZWNJ", "ZWJ", "SHY"},
		},
		{
			name: "WB5: Do not break between most letters",
			text: "안녕 world",
			want: []string{"안녕", "world"},
		},
		{
			name: "WB6,7: Do not break letters across certain punctuation",
			text: "e.g. example.com aujourd'hui",
			want: []string{"e.g", "example.com", "aujourd'hui"},
		},
		{
			name: "WB8,9,10: Do not break on digits",
			text: "88 9원 A10",
			want: []string{"88", "9원", "A10"},
		},
		{
			name: "WB11,12: Do not break between sequences",
			text: "3.14 1,337",
			want: []string{"3.14", "1,337"},
		},
		{
			name: "WB13: Do not break between Katakana",
			text: "ワールド",
			want: []string{"ワールド"},
		},
		{
			name: "WB13a,b: Do not break from extenders",
			text: "NNB\u202fSP Low\u005fLine Under\u203fTie",
			want: []string{"NNB\u202fSP", "Low\u005fLine", "Under\u203fTie"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			subject := tokens.UAX29{}
			got := subject.Tokens(tt.text)

			if !slices.Equal(got, tt.want) {
				t.Errorf("Tokens: want %#v, got %#v", tt.want, got)
			}
		})
	}
}
