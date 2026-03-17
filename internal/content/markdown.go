package content

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
)

type MarkdownStore struct {
	contentPath string
	markdown    goldmark.Markdown
}

func NewMarkdownStore(contentPath string) *MarkdownStore {
	return &MarkdownStore{
		contentPath: contentPath,
		markdown:    goldmark.New(),
	}
}

// Find retrieves a post's filename given its route and locale.
// If no match is found, ("", nil) is returned.
//
// The route is not sanitized, so do that before calling.
func (s *MarkdownStore) Find(route string, locale string) (string, error) {
	filename := filepath.Base(route)
	path := filepath.Join(s.contentPath, filepath.Dir(route))

	// todo: I might want to cache this.
	files, err := scan(path, ".md")
	if err != nil {
		return "", fmt.Errorf("scan %s: %w", path, err)
	}

	if variants, ok := files[filename]; ok {
		return findBestVariant(variants, locale)
	}
	return "", nil
}

func (s *MarkdownStore) Render(w io.Writer, filePath string) error {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", filePath, err)
	}

	err = s.markdown.Convert(contents, w)
	if err != nil {
		return fmt.Errorf("convert %s: %w", filePath, err)
	}

	return nil
}
