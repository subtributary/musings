package markdown

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/subtributary/musings/internal/localization"
	"github.com/yuin/goldmark"
)

type Store struct {
	contentPath string
	markdown    goldmark.Markdown
}

type FileData struct {
	HtmlContent template.HTML
}

func NewStore(contentPath string) *Store {
	return &Store{
		contentPath: contentPath,
		markdown:    goldmark.New(),
	}
}

// Find retrieves a post's filePath given its route and locale
func (s *Store) Find(route string, locale string) (string, error) {
	contentRoot, err := os.OpenRoot(s.contentPath)
	if err != nil {
		return "", fmt.Errorf("could not open content root: %w", err)
	}
	defer func() { _ = contentRoot.Close() }()

	path := filepath.Dir(route)
	files, err := localization.ScanFiles(contentRoot, path, ".md")
	if err != nil {
		return "", fmt.Errorf("could not scan for markdown files: %w", err)
	}

	filename := filepath.Base(route)
	if variants, ok := files[filename]; ok {
		if filePath, err := variants.UseLocale(locale); err == nil {
			return filePath, nil
		}
	}

	return "", fmt.Errorf("could not find locale %q for file %q", locale, filename)
}

func (s *Store) Parse(filePath string) (*FileData, error) {
	contentRoot, err := os.OpenRoot(s.contentPath)
	if err != nil {
		return nil, fmt.Errorf("could not open content root: %w", err)
	}
	defer func() { _ = contentRoot.Close() }()

	contents, err := contentRoot.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", filePath, err)
	}

	// todo: parse metadata at start of file.

	buffer := bytes.Buffer{}
	err = s.markdown.Convert(contents, &buffer)
	if err != nil {
		return nil, fmt.Errorf("convert %q: %w", filePath, err)
	}

	return &FileData{
		HtmlContent: template.HTML(buffer.String()),
	}, nil
}
