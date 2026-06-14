package posts_test

import (
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

	index, err := posts.BuildIndex(files)
	if err != nil {
		t.Fatalf("Error building index: %v", err)
	}

	results := index.List()

	wantIds := []string{"empty", "title", "12", "dec", "2", "1", "jan"}
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

	index, err := posts.BuildIndex(files)
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

func getIds(indexedPosts []posts.IndexedPost) []string {
	var results []string
	for _, p := range indexedPosts {
		results = append(results, p.Path)
	}
	return results
}
