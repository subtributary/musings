package tokens_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/subtributary/musings/internal/tokens"
)

func TestNGram(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		minN, maxN int
		text       string
		want       []string
	}{
		{
			name: "empty",
			minN: 1,
			maxN: 1,
			text: "",
			want: nil,
		},
		{
			name: "unigram",
			minN: 1,
			maxN: 1,
			text: "abcd",
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "bigram",
			minN: 2,
			maxN: 2,
			text: "abcd",
			want: []string{"ab", "bc", "cd"},
		},
		{
			name: "trigram",
			minN: 3,
			maxN: 3,
			text: "abcd",
			want: []string{"abc", "bcd"},
		},
		{
			name: "bigram and trigram",
			minN: 2,
			maxN: 3,
			text: "abcd",
			want: []string{"ab", "bc", "cd", "abc", "bcd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			subject := tokens.NGram{MinN: tt.minN, MaxN: tt.maxN}
			got := subject.Tokens(tt.text)

			if !slices.Equal(got, tt.want) {
				t.Errorf("Tokens = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNGram_JSON(t *testing.T) {
	t.Parallel()

	original := tokens.NGram{
		MinN: 1,
		MaxN: 2,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("error marshaling ngram: %v", err)
	}

	var rebuilt tokens.NGram
	if err := json.Unmarshal(data, &rebuilt); err != nil {
		t.Fatalf("error unmarshaling ngram: %v", err)
	}

	if original != rebuilt {
		t.Errorf("JSON round trip: got %v, want %v", rebuilt, original)
	}
}
