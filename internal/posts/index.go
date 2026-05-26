package posts

import (
	"encoding/json"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/subtributary/search"
	"golang.org/x/text/language"
)

type Index struct {
	wrapped *search.Index
}

func NewIndex() (Index, error) {
	wrapped, err := search.NewIndex(
		search.WithField("title", 5),
		search.WithField("content", 1),
	)
	return Index{wrapped}, err
}

func LoadIndex(dataPath string, locale language.Tag) (Index, error) {
	index := Index{wrapped: &search.Index{}}

	content, err := os.ReadFile(databasePath(dataPath, locale))
	if err != nil {
		return index, fmt.Errorf("read index file: %v", err)
	}

	if err := json.Unmarshal(content, index.wrapped); err != nil {
		return index, fmt.Errorf("unmarshal index: %v", err)
	}

	return index, nil
}

func (idx Index) SaveIndex(dataPath string, locale language.Tag) error {
	content, err := json.Marshal(idx.wrapped)
	if err != nil {
		return fmt.Errorf("marshal index: %v", err)
	}

	filePath := databasePath(dataPath, locale)
	if err := os.WriteFile(filePath, content, os.ModePerm); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}

type SearchResult struct {
	Path  string
	Title string
}

func (idx Index) Search(query string) iter.Seq[SearchResult] {
	return func(yield func(SearchResult) bool) {
		for r := range idx.wrapped.Search(query) {
			if !yield(SearchResult{
				Path:  r.Attachments["att_path"],
				Title: r.Attachments["att_title"],
			}) {
				return
			}
		}
	}
}

func (idx Index) Upsert(path string, post ParsedPost) error {
	id := post.Published.Format("20060102T150405") + path
	path = strings.TrimRight(path, ".md")
	return idx.wrapped.Upsert(id, map[string]string{
		"title":     post.Title,
		"content":   string(post.Content),
		"att_path":  path,
		"att_title": post.Title,
	})
}

func databasePath(dataPath string, locale language.Tag) string {
	if locale == language.Und {
		return filepath.Join(dataPath, "index.json")
	}
	localeStr := strings.ToLower(locale.String())
	return filepath.Join(dataPath, "index."+localeStr+".json")
}
