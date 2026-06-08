package tokens_test

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"

	"github.com/subtributary/musings/internal/tokens"
)

func TestSmart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		text      string
		analyzers map[string][]tokens.Analyzer
		want      []string
	}{
		{
			name: "empty",
			text: "",
			analyzers: map[string][]tokens.Analyzer{
				"Latin": {addPrefix{"Latin"}},
			},
			want: nil,
		},
		{
			name:      "unconfigured",
			text:      "hello world",
			analyzers: map[string][]tokens.Analyzer{},
			want:      nil,
		},
		{
			name:      "empty and unconfigured tokens",
			text:      "",
			analyzers: map[string][]tokens.Analyzer{},
			want:      nil,
		},
		{
			name: "single word",
			text: "hello",
			analyzers: map[string][]tokens.Analyzer{
				"Latin": {addPrefix{"Latin"}},
			},
			want: []string{"Latin:hello"},
		},
		{
			name: "single script",
			text: "hello world",
			analyzers: map[string][]tokens.Analyzer{
				"Latin": {addPrefix{"Latin"}},
			},
			want: []string{"Latin:hello world"},
		},
		{
			name: "two scripts",
			text: "world 안녕",
			analyzers: map[string][]tokens.Analyzer{
				"Latin":  {addPrefix{"Latin"}},
				"Hangul": {addPrefix{"Hangul"}},
			},
			want: []string{"Latin:world ", "Hangul:안녕"},
		},
		{
			name: "leading common",
			text: "--hello world",
			analyzers: map[string][]tokens.Analyzer{
				"Latin": {addPrefix{"Latin"}},
			},
			want: []string{"Latin:hello world"},
		},
		{
			name: "intermixed common",
			text: "hello! 안녕!",
			analyzers: map[string][]tokens.Analyzer{
				"Latin":  {addPrefix{"Latin"}},
				"Hangul": {addPrefix{"Hangul"}},
			},
			want: []string{"Latin:hello! ", "Hangul:안녕!"},
		},
		{
			name: "multiple analyzers",
			text: "hi!",
			analyzers: map[string][]tokens.Analyzer{
				"Latin": {
					addPrefix{"first"},
					addPrefix{"second"},
				},
			},
			want: []string{"second:first:hi!"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			subject := tokens.Smart{}
			for script, analyzers := range tt.analyzers {
				subject.SetAnalyzers(script, analyzers)
			}

			got := subject.Tokens(tt.text)

			// If the tokens are right, the scripts chosen are too.
			if !slices.Equal(got, tt.want) {
				t.Errorf("Tokens returned %#v, want %#v", got, tt.want)
			}
		})
	}
}

type addPrefix struct {
	Prefix string `json:"prefix"`
}

func (a addPrefix) Tokens(text string) []string {
	return []string{a.Prefix + ":" + text}
}

func TestSmart_JSON(t *testing.T) {
	t.Parallel()

	original := &tokens.Smart{}
	original.SetAnalyzers("Arabic", []tokens.Analyzer{
		tokens.Lowercase{},
		tokens.NFKC{},
		tokens.NGram{MinN: 1, MaxN: 2},
		tokens.UAX29{},
	})
	original.SetAnalyzers("Adlam", []tokens.Analyzer{
		tokens.NFKC{},
		tokens.UAX29{},
	})
	original.SetAnalyzers("Ahom", nil)

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Error marshaling: %v", err)
	}

	var rebuilt tokens.Smart
	if err := json.Unmarshal(data, &rebuilt); err != nil {
		t.Fatalf("Error unmarshaling: %v", err)
	}

	wantScripts := original.Scripts()
	slices.Sort(wantScripts)
	gotScripts := rebuilt.Scripts()
	slices.Sort(gotScripts)
	if !slices.Equal(wantScripts, gotScripts) {
		t.Fatalf("Scripts = %#v, want %#v", gotScripts, wantScripts)
	}

	for _, script := range wantScripts {
		wantAnalyzers := original.Analyzers(script)
		gotAnalyzers := rebuilt.Analyzers(script)

		if !slices.EqualFunc(wantAnalyzers, gotAnalyzers, areAnalyzersEqual) {
			t.Errorf("Analyzers(%q) return value is unexpected.", script)
		}
	}
}

func areAnalyzersEqual(a, b tokens.Analyzer) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}
