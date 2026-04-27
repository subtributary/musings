package search

import (
	"cmp"
	"log"
	"math"
	"slices"
)

type BM25 struct {
	k1 float64 // Free parameter. See NewBM25 for details.
	b  float64 // Free parameter. See NewBM25 for details.
}

// NewBM25 creates a new BM25 algorithm context.
//
// k1 is used to approximate the saturation model via a parametric curve.
// With a high value, increases in term frequencies affect scores more.
// With a low value, increases in term frequencies affect scores less.
// For most corpora, a value between 1.2 and 2 is good.
//
// b affects the strength of document-length normalization.
// With a high value, it applies more; with a low value, it applies less.
// For most corpora, a value between 0.5 and 0.8 is good.
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

// Search returns document names sorted by their score for the query.
// The best match is listed first.
// The relative order of equal matches is undefined.
func (bm BM25) Search(corpus Corpus, query []string) []string {
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

// Scores calculates the score of each document for the query.
// A document score means nothing in isolation, but it can be used to compare
// documents for search result ranking. A higher value means a closer match.
func (bm BM25) Scores(corpus Corpus, query []string) map[string]float64 {
	documents := corpus.Documents()

	result := make(map[string]float64, len(documents))
	for _, word := range query {
		idf := bm.idf(corpus, word)
		for name, doc := range documents {
			result[name] += idf * bm.saturation(corpus, doc, word)
		}
	}
	return result
}

// The idf is the importance of the word based on its rarity.
func (bm BM25) idf(corpus Corpus, word string) float64 {
	// A modified Robertson/Sparck Jones formula works well for this.
	docCount := float64(corpus.Size())
	docFrequency := float64(corpus.Count(word))
	return math.Log((docCount-docFrequency+0.5)/(docFrequency+0.5) + 1)
}

// The saturation is how much the word is used in the document.
func (bm BM25) saturation(corpus Corpus, doc Document, word string) float64 {
	// Use the relative document size to normalize the value.
	avgDocSize := corpus.AverageDocumentSize()
	sizeWeight := 1 - bm.b + bm.b*float64(doc.Size())/avgDocSize

	// Return count/size but with normalizations and diminishing returns.
	wordCount := float64(doc.Count(word))
	return wordCount * (bm.k1 + 1) / (wordCount + bm.k1*sizeWeight)
}
