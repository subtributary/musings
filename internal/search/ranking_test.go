package search

import (
	"log"
	"slices"
	"strings"
	"testing"
)

func TestBM25(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		contents []string
		query    string
		rank     []string
	}{
		{
			name:     "empty",
			contents: []string{},
			query:    "test",
			rank:     []string{},
		},
		{
			name:     "single match",
			contents: []string{"blue", "test", "blue tulips"},
			query:    "test",
			rank:     []string{"test", "blue", "blue tulips"},
		},
		{
			name:     "multiple matches",
			contents: []string{"blue", "test", "blue tulips"},
			query:    "blue",
			rank:     []string{"blue", "blue tulips", "test"},
		},
		{
			name:     "overused word",
			contents: []string{"blue", "test", "test test"},
			query:    "test",
			rank:     []string{"test test", "test", "blue"},
		},
		{
			name:     "multiword query",
			contents: []string{"test", "blue", "blue tulips"},
			query:    "tulips blue",
			rank:     []string{"blue tulips", "blue", "test"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			corpus := buildCorpus(t, tt.contents)

			bm := NewBM25(1.2, 0.75)
			query := strings.Split(tt.query, " ")
			rank := bm.Search(corpus, query)

			if !slices.Equal(rank, tt.rank) {
				log.Fatalf("rank: expected %v, got %v", tt.rank, rank)
			}
		})
	}
}

// buildCorpus builds a corpus, naming each document its content.
func buildCorpus(t *testing.T, contents []string) CorpusOld {
	t.Helper()
	corpus := NewCorpusOld()
	for _, content := range contents {
		filename := content
		document := NewDocumentOld(strings.Split(content, " "))
		corpus.Add(filename, document)
	}
	return corpus
}
