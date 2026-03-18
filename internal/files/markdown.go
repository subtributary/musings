package files

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
)

type MarkdownStore struct {
	contentRoot *os.Root
	markdown    goldmark.Markdown
}

type MdFileData struct {
	HtmlContent template.HTML
}

func NewMarkdownStore(contentPath string) (*MarkdownStore, error) {
	contentRoot, err := os.OpenRoot(contentPath)
	if err != nil {
		return nil, fmt.Errorf("could not open content root at %s: %w", contentPath, err)
	}

	return &MarkdownStore{
		contentRoot: contentRoot,
		markdown:    goldmark.New(),
	}, nil
}

func (s *MarkdownStore) Dispose() {
	_ = s.contentRoot.Close()
}

// Find retrieves a post's filename given its route and locale.
// If no match is found, ("", nil) is returned.
func (s *MarkdownStore) Find(route string, locale string) (string, error) {
	filename := filepath.Base(route)
	path := filepath.Dir(route)

	// todo: I might want to cache this.
	files, err := scan(s.contentRoot, path, ".md")
	if err != nil {
		return "", fmt.Errorf("scan %s: %w", path, err)
	}

	if variants, ok := files[filename]; ok {
		return findBestVariant(variants, locale)
	}
	return "", nil
}

// Parse parses the file to extract its metadata and convert its contents to HTML.
func (s *MarkdownStore) Parse(filePath string) (*MdFileData, error) {
	contents, err := s.contentRoot.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", filePath, err)
	}

	// todo: parse metadata at start of file.

	buffer := bytes.Buffer{}
	err = s.markdown.Convert(contents, &buffer)
	if err != nil {
		return nil, fmt.Errorf("convert %s: %w", filePath, err)
	}

	return &MdFileData{
		HtmlContent: template.HTML(buffer.String()),
	}, nil
}
