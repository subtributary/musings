package search

type Document struct {
	tokens []string
}

func NewDocument(tokens []string) Document {
	return Document{tokens: tokens}
}

func (d Document) Size() int {
	return len(d.tokens)
}

// Count returns the number of times the word appears in the document.
func (d Document) Count(word string) int {
	count := 0
	for _, t := range d.tokens {
		if t == word {
			count++
		}
	}
	return count
}

type Corpus struct {
	documents    map[string]Document
	wordCount    int            // Total word count across documents
	docsWithWord map[string]int // Number of documents that contain each word
}

func NewCorpus() Corpus {
	return Corpus{
		documents:    make(map[string]Document),
		docsWithWord: make(map[string]int),
	}
}

/* Actions */

// Add processes and saves a document.
func (c *Corpus) Add(name string, document Document) {
	c.documents[name] = document
	c.wordCount += document.Size()
	for _, word := range removeDuplicates(document.tokens) {
		c.docsWithWord[word]++
	}
}

// Remove removes all data associated with a document.
// It returns a boolean indicating whether the doc was found.
func (c *Corpus) Remove(docName string) bool {
	doc, found := c.documents[docName]
	if found {
		delete(c.documents, docName)
		c.wordCount -= doc.Size()
		for _, word := range removeDuplicates(doc.tokens) {
			c.docsWithWord[word]--
		}
	}
	return found
}

func removeDuplicates(tokens []string) (result []string) {
	visited := make(map[string]struct{}, len(tokens))
	for _, token := range tokens {
		if _, ok := visited[token]; !ok {
			result = append(result, token)
			visited[token] = struct{}{}
		}
	}
	return result
}

/* Queries */

func (c *Corpus) AverageDocumentSize() float64 {
	if docCount := len(c.documents); docCount > 0 {
		return float64(c.wordCount) / float64(docCount)
	}
	return 0
}

// Count returns the number of documents containing a word.
func (c *Corpus) Count(word string) int {
	return c.docsWithWord[word]
}

// Documents returns a map of document names to documents.
func (c *Corpus) Documents() map[string]Document {
	return c.documents
}

// Size returns the number of documents in the corpus.
func (c *Corpus) Size() int {
	return len(c.documents)
}
