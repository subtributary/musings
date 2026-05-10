package search_test

import (
	"strings"
	"testing"

	"github.com/subtributary/musings/internal/search"
)

func TestDocument_Count(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content map[search.Field]string
		word    string
		want    map[search.Field]int
	}{
		{
			name:    "no streams",
			content: map[search.Field]string{"f1": ""},
			word:    "apple",
			want:    map[search.Field]int{"f1": 0},
		},
		{
			name:    "word does not exist in any stream",
			content: map[search.Field]string{"f1": "apple banana cucumber"},
			word:    "date",
			want:    map[search.Field]int{"f1": 0},
		},
		{
			name: "word exists in one stream",
			content: map[search.Field]string{
				"f1": "apple banana cucumber",
				"f2": "zebra yak xolo",
			},
			word: "cucumber",
			want: map[search.Field]int{"f1": 1, "f2": 0},
		},
		{
			name: "word exists in two streams",
			content: map[search.Field]string{
				"f1": "zucchini yogurt white rice",
				"f2": "zebra yak xolo white tiger",
			},
			word: "white",
			want: map[search.Field]int{"f1": 1, "f2": 1},
		},
		{
			name: "word exists twice in a stream",
			content: map[search.Field]string{
				"f1": "eggplant fish galjoen fish",
			},
			word: "fish",
			want: map[search.Field]int{"f1": 2},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			document := search.NewDocument()
			for field, content := range tt.content {
				tokens := strings.Split(content, " ")
				if tokens[0] == "" {
					tokens = []string{}
				}
				document.SetStream(field, tokens)
			}

			for field, want := range tt.want {
				got := document.Count(field, tt.word)
				if got != want {
					t.Errorf("document.Count(%q, %q): got %d, want %d", field, tt.word, got, want)
				}
			}
		})
	}
}

func TestDocument_Length(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content map[search.Field]string
		want    map[search.Field]int
	}{
		{
			name:    "empty stream",
			content: map[search.Field]string{"f1": ""},
			want:    map[search.Field]int{"f1": 0},
		},
		{
			name: "streams of different lengths",
			content: map[search.Field]string{
				"f1": "apple banana cucumber",
				"f2": "whale",
			},
			want: map[search.Field]int{"f1": 3, "f2": 1},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			document := search.NewDocument()
			for field, content := range tt.content {
				tokens := strings.Split(content, " ")
				if tokens[0] == "" {
					tokens = []string{}
				}
				document.SetStream(field, tokens)
			}

			for field, want := range tt.want {
				got := document.Length(field)
				if got != want {
					t.Errorf("document.Length(%q): got %d, want %d", field, got, want)
				}
			}
		})
	}
}

func TestCorpus(t *testing.T) {
	t.Parallel()
	corpus := search.NewCorpus()

	assertCount := func(term string, want int) {
		t.Helper()
		if got := corpus.Count(term); got != want {
			t.Errorf("corpus.Count(%q): got %d, want %d", term, got, want)
		}
	}

	assertSize := func(want int) {
		t.Helper()
		if got := corpus.Size(); got != want {
			t.Errorf("corpus.Size(): got %d, want %d", got, want)
		}
	}

	assertAvgDocSize := func(want float64) {
		t.Helper()
		if got := corpus.AvgStreamLength(""); got != want {
			t.Errorf("corpus.AvgStreamLength(): got %v, want %v", got, want)
		}
	}

	// Populate corpus
	corpus.Upsert("one", createDocument("one"))
	corpus.Upsert("two", createDocument("two two"))
	corpus.Upsert("three", createDocument("three three three"))
	assertAvgDocSize(2)
	assertCount("three", 1)
	assertSize(3)
	if t.Failed() {
		t.FailNow()
	}

	// Replace existing document
	corpus.Upsert("one", createDocument("one two three four"))
	assertAvgDocSize(3)
	assertCount("three", 2)
	assertSize(3)
	if t.Failed() {
		t.FailNow()
	}

	// Remove document
	corpus.Remove("one")
	assertAvgDocSize(2.5)
	assertCount("one", 0)
	assertSize(2)
	if t.Failed() {
		t.FailNow()
	}

	// Remove remaining documents
	corpus.Remove("two")
	corpus.Remove("three")
	assertAvgDocSize(0)
	assertCount("two", 0)
	assertSize(0)
	if t.Failed() {
		t.FailNow()
	}

	// Remove nonexistent
	corpus.Remove("missing")
	assertAvgDocSize(0.0)
	assertSize(0)
}

func createDocument(text string) search.Document {
	tokens := strings.Split(text, " ")
	document := search.NewDocument()
	document.SetStream("", tokens)
	return document
}
