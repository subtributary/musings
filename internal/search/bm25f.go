package search

import (
	"cmp"
	"math"
	"slices"
)

type StreamConfig struct {
	Weight float64

	// B is the strength of length normalization.
	// For 0, no normalization is performed.
	// For 1, results are scaled to the average document length.
	B float64
}

type BM25F struct {
	corpus        Corpus
	streamConfigs []StreamConfig
	k1            float64
}

// Rank returns document names sorted by their score for the query.
// The best match is listed first.
// Equal matches are sorted lexigraphically by filename.
func (bm BM25F) Rank(query []string) []string {
	scores := bm.Scores(query)

	names := make([]string, len(scores))
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

func (bm BM25F) Scores(query []string) map[string]float64 {
	result := make(map[string]float64, bm.corpus.Size())
	for _, doc := range bm.corpus.Documents() {
		result[doc.Name] = bm.score(doc, query)
	}
	return result
}

func (bm BM25F) score(doc Document, query []string) (result float64) {
	for _, term := range query {
		// The score is the product of the saturation and the importance (IDF).
		termFreq := bm.termFrequency(doc, term)
		saturation := termFreq / (bm.k1 + termFreq)
		result += saturation * bm.idf(term)
	}
	return
}

// idf returns the relative importance of a word based on its rarity.
func (bm BM25F) idf(term string) float64 {
	// For the IDF, we apply a modified Robertson/Sparck Jones formula across
	// all streams. There are rare scenarios where this does not yield good
	// results. We will ignore the problem until it shows itself in practice.
	docCount := float64(bm.corpus.Size())
	docFreq := float64(bm.corpus.Count(term))
	return math.Log((docCount-docFreq+0.5)/(docFreq+0.5) + 1)
}

// termFrequency returns the normalized weighted frequency of a term within the
// document across all stream.
func (bm BM25F) termFrequency(doc Document, term string) (result float64) {
	for i, config := range bm.streamConfigs {
		avgStreamLen := bm.corpus.AverageStreamLength(i)
		if avgStreamLen == 0.0 {
			continue
		}

		// Normalize results when the stream length is far from average.
		streamLen := float64(doc.Length(i))
		lengthNorm := 1 - config.B + config.B*streamLen/avgStreamLen

		// Simple weighted summation with normalization.
		termFreq := float64(doc.Count(i, term))
		result += config.Weight * termFreq / lengthNorm
	}
	return
}
