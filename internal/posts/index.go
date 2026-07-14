package posts

import (
	"cmp"
	"encoding/json"
	"fmt"
	"iter"
	"slices"
	"strings"
	"time"

	"github.com/computerghost/bm25f"
	"github.com/subtributary/musings/internal/tokens"
)

const (
	fieldTitle   = "title"
	fieldContent = "content"

	metadataBylines   = "bylines"
	metadataPath      = "path"
	metadataPublished = "published"
	metadataSummary   = "summary"
	metadataThumbnail = "thumbnail"
	metadataTitle     = "title"
)

type IndexedPost struct {
	Path      string
	Bylines   []string
	Published time.Time
	Summary   string
	Thumbnail string
	Title     string
}

type Index struct {
	ranker    *bm25f.BM25F
	corpus    bm25f.Corpus
	tokenizer tokens.Analyzer
}

// NewIndex creates an empty index ready to be populated.
// Once the initial population is done, call MakeThreadSafe.
func NewIndex() *Index {
	ranker := bm25f.New()
	// Errors are ignored because they only happen for weight < 0.
	_ = ranker.SetWeight(fieldTitle, 5.0)
	_ = ranker.SetWeight(fieldContent, 1.0)

	return &Index{
		ranker:    ranker,
		corpus:    bm25f.NewSimpleCorpus(),
		tokenizer: tokens.NewDefaultTokenizer(),
	}
}

// List lists all posts in descending order of publication date.
// Posts with publication dates in the future are omitted.
func (idx *Index) List() iter.Seq[*IndexedPost] {
	var results []*IndexedPost
	for _, doc := range idx.corpus.Documents() {
		results = append(results, docToPost(doc))
	}

	slices.SortFunc(results, func(a, b *IndexedPost) int {
		// compareTimes parameters are reversed to sort descending.
		if c := b.Published.Compare(a.Published); c != 0 {
			return c
		}
		if c := cmp.Compare(a.Title, b.Title); c != 0 {
			return c
		}
		return cmp.Compare(a.Path, b.Path)
	})

	now := time.Now()

	return func(yield func(p *IndexedPost) bool) {
		for _, post := range results {
			if post.Published.After(now) {
				continue
			}
			if !yield(post) {
				return
			}
		}
	}
}

// MakeConcurrent makes the index thread-safe at the cost of update speed.
// Call this after the index's initial population.
func (idx *Index) MakeConcurrent() error {
	sc, ok := idx.corpus.(*bm25f.SimpleCorpus)
	if !ok {
		return fmt.Errorf("index corpus is already thread-safe")
	}

	idx.corpus = bm25f.NewSyncCorpus(sc)

	return nil
}

// Remove removes a post from the index.
func (idx *Index) Remove(name string) {
	idx.corpus.Remove(name)
}

// Search returns the posts matching the query,
// sorted by match score with the best match first.
// Posts with publication dates in the future are omitted.
func (idx *Index) Search(query string) iter.Seq[*IndexedPost] {
	if strings.TrimSpace(query) == "" {
		return idx.List()
	}

	queryToks := idx.tokenizer.Tokens(query)
	scores := idx.ranker.Score(idx.corpus, queryToks)

	scores = slices.DeleteFunc(scores, func(r bm25f.Result) bool {
		return r.Score == 0
	})

	slices.SortFunc(scores, func(a, b bm25f.Result) int {
		if c := cmp.Compare(b.Score, a.Score); c != 0 {
			return c
		}
		return cmp.Compare(a.ID, b.ID)
	})

	now := time.Now()

	return func(yield func(p *IndexedPost) bool) {
		for _, result := range scores {
			post := docToPost(result.Document)
			if post.Published.After(now) {
				continue
			}
			if !yield(post) {
				return
			}
		}
	}
}

// Upsert adds a post to the index.
func (idx *Index) Upsert(name string, post *ParsedPost) {
	routePath := strings.TrimSuffix(name, ".md")

	bylines, _ := json.Marshal(post.Bylines)
	content := idx.tokenizer.Tokens(string(post.Content))
	title := idx.tokenizer.Tokens(post.Title)

	var published string
	if !post.Published.IsZero() {
		published = post.Published.Format(time.RFC3339)
	}

	idx.corpus.Upsert(routePath, bm25f.NewDocument(
		bm25f.WithField(fieldTitle, title),
		bm25f.WithField(fieldContent, content),
		bm25f.WithMetadata(metadataBylines, string(bylines)),
		bm25f.WithMetadata(metadataPath, routePath),
		bm25f.WithMetadata(metadataPublished, published),
		bm25f.WithMetadata(metadataSummary, post.Summary),
		bm25f.WithMetadata(metadataThumbnail, post.Thumbnail),
		bm25f.WithMetadata(metadataTitle, post.Title),
	))
}

// docToPost converts a bm25f.Document to an IndexedPost.
// Missing or malformed metadata is degraded to zero values.
func docToPost(doc *bm25f.Document) *IndexedPost {
	p := &IndexedPost{}

	p.Path, _ = doc.Metadata(metadataPath)
	p.Summary, _ = doc.Metadata(metadataSummary)
	p.Thumbnail, _ = doc.Metadata(metadataThumbnail)
	p.Title, _ = doc.Metadata(metadataTitle)

	bylinesStr, _ := doc.Metadata(metadataBylines)
	_ = json.Unmarshal([]byte(bylinesStr), &p.Bylines)

	publishedStr, _ := doc.Metadata(metadataPublished)
	p.Published, _ = time.Parse(time.RFC3339, publishedStr)

	return p
}
