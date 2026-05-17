package posts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/subtributary/search"
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

func LoadIndex(dataPath string, locale string) (Index, error) {
	index := Index{wrapped: &search.Index{}}
	filePath := databasePath(dataPath, locale)
	if content, err := os.ReadFile(filePath); err != nil {
		return index, fmt.Errorf("read index file: %v", err)
	} else if err = json.Unmarshal(content, index.wrapped); err != nil {
		return index, fmt.Errorf("unmarshal index: %v", err)
	}
	return index, nil
}

func (idx Index) SaveIndex(dataPath string, locale string) error {
	filePath := databasePath(dataPath, locale)
	if content, err := json.Marshal(idx.wrapped); err != nil {
		return fmt.Errorf("marshal index: %v", err)
	} else if err := os.WriteFile(filePath, content, os.ModePerm); err != nil {
		return fmt.Errorf("write file: %v", err)
	}
	return nil
}

func (idx Index) Upsert(path string, post ParsedPost) error {
	return idx.wrapped.Upsert(path, map[string]string{
		"title":         post.Title,
		"content":       post.Content,
		"display_title": post.Title,
	})
}

func databasePath(dataPath, locale string) string {
	if locale == "" {
		return filepath.Join(dataPath, "index.json")
	} else {
		return filepath.Join(dataPath, "index."+locale+".json")
	}
}
