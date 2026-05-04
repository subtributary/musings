package search

type Index struct {
	corpus CorpusOld
	ranker BM25
}

func NewIndex() Index {
	return Index{
		corpus: NewCorpusOld(),
		ranker: NewBM25(1.2, 0.75),
	}
}

// Add analyzes a document and adds its data to the corpus.
func (idx Index) Add(name string, text string) {
	tokens := Tokenize(text)
	idx.corpus.Add(name, NewDocumentOld(tokens))
}

func (idx Index) Remove(name string) {
	idx.corpus.Remove(name)
}

// Search searches the corpus for the text, and it returns the names
// of the documents sorted with the best matches first.
func (idx Index) Search(query string) []string {
	queryTokens := Tokenize(query)
	return idx.ranker.Search(idx.corpus, queryTokens)
}
