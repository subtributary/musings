package search

import (
	"maps"
	"slices"
)

type Stream []string

type Document struct {
	Name    string
	streams []Stream
}

// Count returns the number of times a term appears in a stream.
func (d Document) Count(i int, term string) (result int) {
	for _, t := range d.streams[i] {
		if t == term {
			result++
		}
	}
	return
}

// Length returns the number of tokens in a stream.
func (d Document) Length(i int) int {
	return len(d.streams[i])
}

func (d Document) UniqueWords() (result []string) {
	visited := make(map[string]struct{})
	for _, stream := range d.streams {
		for _, token := range stream {
			if _, ok := visited[token]; !ok {
				result = append(result, token)
				visited[token] = struct{}{}
			}
		}
	}
	return
}

type Corpus struct {
	documents    map[string]Document
	totalLengths []int          // Total stream lengths
	docsWithTerm map[string]int // Number of documents containing each term
}

// Add processes and saves a document.
func (c *Corpus) Add(name string, document Document) {
	c.Remove(name)

	c.documents[name] = document
	for i, stream := range document.streams {
		c.totalLengths[i] += len(stream)
	}
	for _, word := range document.UniqueWords() {
		c.docsWithTerm[word]++
	}
}

// Remove removes all data associated with a document.
func (c *Corpus) Remove(name string) {
	if doc, ok := c.documents[name]; ok {
		delete(c.documents, name)
		for i, stream := range doc.streams {
			c.totalLengths[i] -= len(stream)
		}
		for _, word := range doc.UniqueWords() {
			c.docsWithTerm[word]--
		}
	}
}

// AverageStreamLength returns the average length of a stream across the corpus.
func (c *Corpus) AverageStreamLength(i int) float64 {
	if docCount := len(c.documents); docCount > 0 {
		return float64(c.totalLengths[i]) / float64(docCount)
	}
	return 0
}

// Count returns the number of documents that contain a term.
func (c *Corpus) Count(term string) int {
	return c.docsWithTerm[term]
}

// Documents returns the documents in the corpus.
func (c *Corpus) Documents() []Document {
	return slices.Collect(maps.Values(c.documents))
}

// Size returns the number of documents in the corpus.
func (c *Corpus) Size() int {
	return len(c.documents)
}
