package search_test

import (
	"log"
	"slices"
	"strings"
	"testing"

	"github.com/subtributary/musings/internal/search"
)

func TestBM25F(t *testing.T) {
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

			corpus := search.NewCorpus()
			for _, content := range tt.contents {
				filename := content
				document := search.NewDocument()
				document.SetStream("", strings.Split(content, " "))
				corpus.Upsert(filename, document)
			}

			fieldConfigs := []search.FieldConfig{
				{
					Name:   "",
					Weight: 1.0,
					B:      0.72,
				},
			}
			bm := search.NewBM25F(fieldConfigs)
			query := strings.Split(tt.query, " ")
			rank := bm.Rank(corpus, query)

			if !slices.Equal(rank, tt.rank) {
				log.Fatalf("rank: expected %v, got %v", tt.rank, rank)
			}
		})
	}
}
