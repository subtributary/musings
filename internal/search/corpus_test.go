package search

import (
	"strings"
	"testing"
)

func TestDocument_Count(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		word    string
		want    int
	}{
		{
			name:    "word exists in document",
			content: "one two two three three three",
			word:    "two",
			want:    2,
		},
		{
			name:    "word does not exist in document",
			content: "one two two three three three",
			word:    "zero",
			want:    0,
		},
		{
			name:    "document is empty",
			content: "",
			word:    "zero",
			want:    0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			doc := NewDocumentOld(strings.Split(tt.content, " "))
			freq := doc.Count(tt.word)
			if freq != tt.want {
				t.Errorf("Frequency: got %v, want %v", freq, tt.want)
			}
		})
	}
}

func TestCorpus(t *testing.T) {
	t.Parallel()
	corpus := NewCorpusOld()

	assertCount := func(word string, want int) {
		t.Helper()
		if got := corpus.Count(word); got != want {
			t.Errorf("Corpus.Count(%q): got %d, want %d", word, got, want)
		}
	}

	assertSize := func(want int) {
		t.Helper()
		if got := corpus.Size(); got != want {
			t.Errorf("Corpus.Size(): got %d, want %d", got, want)
		}
	}

	assertAvgDocSize := func(want float64) {
		t.Helper()
		if got := corpus.AverageDocumentSize(); got != want {
			t.Errorf("Corpus.AverageDocumentSize(): got %f, want %f", got, want)
		}
	}

	// Populate corpus
	corpus.Add("one", NewDocumentOld([]string{"one"}))
	corpus.Add("two", NewDocumentOld([]string{"two", "two"}))
	corpus.Add("three", NewDocumentOld([]string{"three", "three", "three"}))
	assertAvgDocSize(2.0)
	assertCount("three", 1)
	assertSize(3)
	if t.Failed() {
		t.FailNow()
	}

	// Remove document
	if !corpus.Remove("one") {
		t.Fatalf("Corpus.Remove(): got false, want true")
	}
	assertAvgDocSize(2.5)
	assertCount("one", 0)
	assertSize(2)
	if t.Failed() {
		t.FailNow()
	}

	// Remove remaining documents
	if !corpus.Remove("two") {
		t.Fatalf("Corpus.Remove(): got false, want true")
	}
	if !corpus.Remove("three") {
		t.Fatalf("Corpus.Remove(): got false, want true")
	}
	assertAvgDocSize(0.0)
	assertCount("two", 0)
	assertSize(0)
	if t.Failed() {
		t.FailNow()
	}

	// Remove nonexistent
	if corpus.Remove("missing") {
		t.Fatalf("Corpus.Remove(): got true, want false")
	}
	assertAvgDocSize(0.0)
	assertSize(0)
}
