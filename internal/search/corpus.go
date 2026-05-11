package search

type Stream struct {
	length     int
	termCounts map[string]int
}

func newStream(tokens []string) (s Stream) {
	s.length = len(tokens)
	s.termCounts = make(map[string]int)
	for _, term := range tokens {
		s.termCounts[term]++
	}
	return
}

type Document struct {
	streams map[Field]Stream
}

func NewDocument() (d Document) {
	d.streams = make(map[Field]Stream)
	return
}

func (d Document) SetStream(field Field, tokens []string) {
	d.streams[field] = newStream(tokens)
}

// Count returns the number of times a term appears in a stream.
func (d Document) Count(field Field, term string) (result int) {
	return d.streams[field].termCounts[term]
}

// Length returns the number of tokens in a stream.
func (d Document) Length(field Field) int {
	return d.streams[field].length
}

func (d Document) uniqueWords() (result []string) {
	visited := make(map[string]struct{})
	for _, stream := range d.streams {
		for term := range stream.termCounts {
			result = append(result, term)
			visited[term] = struct{}{}
		}
	}
	return
}

type Corpus struct {
	documents    map[string]Document
	docsWithTerm map[string]int // Number of documents containing each term
	totalLengths map[Field]int  // Total stream lengths
}

func NewCorpus() Corpus {
	return Corpus{
		documents:    make(map[string]Document),
		docsWithTerm: make(map[string]int),
		totalLengths: make(map[Field]int),
	}
}

// Upsert processes and saves a document.
func (c *Corpus) Upsert(name string, document Document) {
	c.Remove(name)

	c.documents[name] = document
	for field, stream := range document.streams {
		c.totalLengths[field] += stream.length
	}
	for _, word := range document.uniqueWords() {
		c.docsWithTerm[word]++
	}
}

// Remove removes all data associated with a document.
func (c *Corpus) Remove(name string) {
	if doc, ok := c.documents[name]; ok {
		delete(c.documents, name)
		for field, stream := range doc.streams {
			c.totalLengths[field] -= stream.length
		}
		for _, word := range doc.uniqueWords() {
			c.docsWithTerm[word]--
			if c.docsWithTerm[word] == 0 {
				delete(c.docsWithTerm, word)
			}
		}
	}
}

// AvgStreamLength returns the average length of a stream across the corpus.
func (c *Corpus) AvgStreamLength(field Field) float64 {
	if docCount := len(c.documents); docCount > 0 {
		return float64(c.totalLengths[field]) / float64(docCount)
	}
	return 0
}

// Count returns the number of documents that contain a term.
func (c *Corpus) Count(term string) int {
	return c.docsWithTerm[term]
}

// Documents returns the documents in the corpus.
func (c *Corpus) Documents() map[string]Document {
	return c.documents
}

// Size returns the number of documents in the corpus.
func (c *Corpus) Size() int {
	return len(c.documents)
}
