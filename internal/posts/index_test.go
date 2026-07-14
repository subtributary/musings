package posts_test

import (
	"fmt"
	"io/fs"
	"path"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/subtributary/musings/internal/posts"
)

func TestIndex_List(t *testing.T) {
	t.Parallel()

	files := fstest.MapFS{
		"1.md":      {Data: []byte("---\npublished: 2025-01-01\n---\n# 1")},
		"jan.md":    {Data: []byte("---\npublished: 2025-01-01\n---\n# Jan")},
		"2.md":      {Data: []byte("---\npublished: 2025-02-01\n---\n# Feb")},
		"12.md":     {Data: []byte("---\npublished: 2025-12-01\n---\n# Same")},
		"dec.md":    {Data: []byte("---\npublished: 2025-12-01\n---\n# Same")},
		"future.md": {Data: []byte("---\npublished: 3000-01-01\n---\n# 3000")},
		"title.md":  {Data: []byte("#No Date")},
		"empty.md":  {Data: []byte("")},
		"ignored":   {Data: []byte("")},
	}

	index, err := buildIndex(files)
	if err != nil {
		t.Fatalf("Error building index: %v", err)
	}

	results := index.List()

	wantIds := []string{"12", "dec", "2", "1", "jan", "empty", "title"}
	gotIds := getIds(slices.Collect(results))
	if !slices.Equal(wantIds, gotIds) {
		t.Errorf("List() = %v, want %v", gotIds, wantIds)
	}
}

func TestIndex_Search(t *testing.T) {
	t.Parallel()

	files := fstest.MapFS{
		"jan 1.md":  {Data: []byte("---\npublished: 2025-01-01\n---\n# Jan 1")},
		"jan 2.md":  {Data: []byte("---\npublished: 2025-01-02\n---\n# Jan 2")},
		"feb 1.md":  {Data: []byte("---\npublished: 2025-02-01\n---\n# Feb 1")},
		"feb 2.md":  {Data: []byte("---\npublished: 2025-02-02\n---\n# Feb 2")},
		"future.md": {Data: []byte("---\npublished: 3000-01-01\n---\n# Jan 3000")},
		"ignored":   {Data: []byte("---\npublished: 2025-01-03\n---\n# Jan 3")},
	}

	index, err := buildIndex(files)
	if err != nil {
		t.Fatalf("Error building index: %v", err)
	}

	results := index.Search("jan")

	wantIds := []string{"jan 1", "jan 2"}
	gotIds := getIds(slices.Collect(results))
	if !slices.Equal(wantIds, gotIds) {
		t.Errorf("List() = %v, want %v", gotIds, wantIds)
	}
}

func buildIndex(contentFS fs.FS) (*posts.Index, error) {
	index := posts.NewIndex()

	versionURL := func(_, name string) string { return name }
	parser := posts.NewParser(versionURL)

	err := fs.WalkDir(contentFS, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.Type().IsRegular() || path.Ext(d.Name()) != ".md" {
			return nil
		}

		post, err := parser.ParseFile(contentFS, filePath)
		if err != nil {
			return fmt.Errorf("parse post %q: %w", filePath, err)
		}

		index.Upsert(filePath, post)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return index, nil
}

func getIds(indexedPosts []*posts.IndexedPost) []string {
	var results []string
	for _, p := range indexedPosts {
		results = append(results, p.Path)
	}
	return results
}
