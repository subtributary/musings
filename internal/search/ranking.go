package search

import (
	"cmp"
	"log"
	"slices"
)

type BM25 struct {
	k1 float64 // Free parameter. See NewBM25 for details.
	b  float64 // Free parameter. See NewBM25 for details.
}

// NewBM25 creates a new BM25 algorithm context.
//
// Reference: https://www.staff.city.ac.uk/~sbrp622/papers/foundations_bm25_review.pdf
func NewBM25(k1 float64, b float64) BM25 {
	if k1 < 0 || b < 0 || b > 1 {
		log.Fatalf("NewBM25(): k1 %f, b %f", k1, b)
	}
	return BM25{
		k1: k1,
		b:  b,
	}
}

// Scores calculates the score of each document for the query.
// A document score means nothing in isolation, but it can be used to compare
// documents for search result ranking. A higher value means a closer match.
func (bm BM25) Scores(corpus CorpusOld, query []string) map[string]float64 {
	documents := corpus.Documents()

	result := make(map[string]float64, len(documents))
	for _, word := range query {
		idf := corpus.IDF(word)
		for name, doc := range documents {
			result[name] += idf * bm.saturation(corpus, doc, word)
		}
	}
	return result
}

// The saturation is how much the word is used in the document.
func (bm BM25) saturation(corpus CorpusOld, doc DocumentOld, word string) float64 {
	// Use the relative document size to normalize the value.
	avgDocSize := corpus.AverageDocumentSize()
	sizeWeight := 1 - bm.b + bm.b*float64(doc.Size())/avgDocSize

	// Return count/size but with normalizations and diminishing returns.
	wordCount := float64(doc.Count(word))
	return wordCount * (bm.k1 + 1) / (wordCount + bm.k1*sizeWeight)
}

// Search returns document names sorted by their score for the query.
// The best match is listed first.
// The relative order of equal matches is undefined.
func (bm BM25) Search(corpus CorpusOld, query []string) []string {
	scores := bm.Scores(corpus, query)

	names := make([]string, 0, len(scores))
	for name := range scores {
		names = append(names, name)
	}

	slices.SortFunc(names, func(a, b string) int {
		if c := cmp.Compare(scores[b], scores[a]); c != 0 {
			return c
		}
		return cmp.Compare(a, b) // Deterministic fallback
	})

	return names
}
