package posts

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"iter"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/computerghost/bm25f"
	"github.com/subtributary/musings/internal/tokens"
)

const (
	fieldTitle   = "title"
	fieldContent = "content"

	metadataTitle     = "title"
	metadataPath      = "path"
	metadataPublished = "published"
	metadataBylines   = "bylines"
)

type IndexedPost struct {
	Path      string
	Bylines   []string
	Published *time.Time
	Title     string
}

func (p IndexedPost) isPublished(now time.Time) bool {
	return p.Published == nil || !p.Published.After(now)
}

type Index struct {
	ranker    *bm25f.BM25F
	corpus    *bm25f.Corpus
	tokenizer tokens.Analyzer
}

// NewIndex creates an empty index ready to be populated.
func NewIndex() *Index {
	s := &Index{
		ranker:    bm25f.New(),
		corpus:    bm25f.NewCorpus(),
		tokenizer: tokens.NewDefaultTokenizer(),
	}

	// Errors are ignored because they only happen for weight < 0.
	_ = s.ranker.SetWeight(fieldTitle, 5.0)
	_ = s.ranker.SetWeight(fieldContent, 1.0)

	return s
}

// BuildIndex builds an entirely new index.
func BuildIndex(ctx context.Context, contentRoot fs.FS) (*Index, error) {
	index := NewIndex()

	parser := NewParser()

	err := fs.WalkDir(contentRoot, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if !d.Type().IsRegular() || path.Ext(d.Name()) != ".md" {
			return nil
		}

		post, err := parser.ParseFile(contentRoot, filePath)
		if err != nil {
			return fmt.Errorf("parse post %q: %w", filePath, err)
		}

		index.upsert(filePath, post)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return index, nil
}

// List lists all posts in order of publication with the most recent one first.
// Posts with publication dates in the future are omitted.
func (s *Index) List() iter.Seq[IndexedPost] {
	var results []IndexedPost
	for _, doc := range s.corpus.Documents() {
		results = append(results, docToPost(doc))
	}

	slices.SortFunc(results, func(a, b IndexedPost) int {
		// compareTimes parameters are reversed to sort descending.
		// This also causes nil times to be sorted before non-nil times.
		if c := compareTimes(b.Published, a.Published); c != 0 {
			return c
		}
		if c := cmp.Compare(a.Title, b.Title); c != 0 {
			return c
		}
		return cmp.Compare(a.Path, b.Path)
	})

	now := time.Now()

	return func(yield func(p IndexedPost) bool) {
		for _, post := range results {
			if !post.isPublished(now) {
				continue
			}
			if !yield(post) {
				return
			}
		}
	}
}

// Search returns the posts matching the query,
// sorted by match score with the best match first.
// Posts with publication dates in the future are omitted.
func (s *Index) Search(query string) iter.Seq[IndexedPost] {
	if strings.TrimSpace(query) == "" {
		return s.List()
	}

	queryToks := s.tokenizer.Tokens(query)
	scores := s.ranker.Score(s.corpus, queryToks)

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

	return func(yield func(p IndexedPost) bool) {
		for _, result := range scores {
			post := docToPost(result.Document)
			if !post.isPublished(now) {
				continue
			}
			if !yield(post) {
				return
			}
		}
	}
}

// upsert adds a post to the index.
func (s *Index) upsert(path string, post ParsedPost) {
	routePath := strings.TrimSuffix(path, ".md")
	title := s.tokenizer.Tokens(post.Title)
	content := s.tokenizer.Tokens(string(post.Content))

	var published string
	if post.Published != nil {
		published = post.Published.Format(time.RFC3339)
	}

	bylines, _ := json.Marshal(post.Bylines)

	s.corpus.Upsert(routePath, bm25f.NewDocument(
		bm25f.WithField(fieldTitle, title),
		bm25f.WithField(fieldContent, content),
		bm25f.WithMetadata(metadataTitle, post.Title),
		bm25f.WithMetadata(metadataPath, routePath),
		bm25f.WithMetadata(metadataPublished, published),
		bm25f.WithMetadata(metadataBylines, string(bylines)),
	))
}

// docToPost converts a bm25f.Document to an IndexedPost.
// Missing or malformed metadata is degrated to zero values.
func docToPost(doc *bm25f.Document) (p IndexedPost) {
	p.Title, _ = doc.Metadata(metadataTitle)
	p.Path, _ = doc.Metadata(metadataPath)

	publishedStr, _ := doc.Metadata(metadataPublished)
	if published, err := time.Parse(time.RFC3339, publishedStr); err == nil {
		p.Published = &published
	}

	bylinesStr, _ := doc.Metadata(metadataBylines)
	_ = json.Unmarshal([]byte(bylinesStr), &p.Bylines)

	return p
}

// compareTimes is like time.Time.Compare but handles nil times.
// Nil times are treated as occurring after non-Nil times.
func compareTimes(a *time.Time, b *time.Time) int {
	switch {
	case a == nil && b == nil:
		return 0
	case a == nil:
		return 1
	case b == nil:
		return -1
	default:
		return a.Compare(*b)
	}
}
