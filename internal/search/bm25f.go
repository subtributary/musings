package search

import (
	"cmp"
	"math"
	"slices"
)

type Field string

type FieldConfig struct {
	Name   Field
	Weight float64

	// B is the strength of length normalization.
	// For 0, no normalization is performed.
	// For 1, results are scaled to the average document length.
	B float64
}

type BM25F struct {
	fieldConfigs []FieldConfig
	k1           float64
}

func NewBM25F(fieldConfigs []FieldConfig) (bm BM25F) {
	bm.fieldConfigs = fieldConfigs
	bm.k1 = 1.2
	return
}

// Rank returns document names sorted by their score for the query.
// The best match is listed first.
// Equal matches are sorted lexigraphically by filename.
func (bm BM25F) Rank(corpus Corpus, query []string) []string {
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

func (bm BM25F) Scores(corpus Corpus, query []string) map[string]float64 {
	// todo: remove duplicates in query
	result := make(map[string]float64, corpus.Size())
	for _, term := range query {
		idf := bm.idf(corpus, term)
		for name, doc := range corpus.Documents() {
			// The score is the saturation multiplied by the importance (IDF).
			termFreq := bm.termFrequency(corpus, doc, term)
			saturation := termFreq / (bm.k1 + termFreq)
			result[name] += saturation * idf
		}
	}
	return result
}

// idf returns the relative importance of a word based on its rarity.
func (bm BM25F) idf(corpus Corpus, term string) float64 {
	// For the IDF, we apply a modified Robertson/Sparck Jones formula across
	// all streams. There are rare scenarios where this does not yield good
	// results. We will ignore the problem until it shows itself in practice.
	docCount := float64(corpus.Size())
	docFreq := float64(corpus.Count(term))
	return math.Log((docCount-docFreq+0.5)/(docFreq+0.5) + 1)
}

// termFrequency returns the normalized weighted frequency of a term within the
// document across all stream.
func (bm BM25F) termFrequency(corpus Corpus, doc Document, term string) (result float64) {
	for _, config := range bm.fieldConfigs {
		avgStreamLen := corpus.AvgStreamLength(config.Name)
		if avgStreamLen == 0.0 {
			continue
		}

		// Normalize results when the stream length is far from average.
		streamLen := float64(doc.Length(config.Name))
		lengthNorm := 1 - config.B + config.B*streamLen/avgStreamLen

		// Simple weighted summation with normalization.
		termFreq := float64(doc.Count(config.Name, term))
		result += config.Weight * termFreq / lengthNorm
	}
	return
}
