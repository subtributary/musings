package posts

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/subtributary/musings/internal/localization"
	"github.com/yuin/goldmark"
)

type Store struct {
	rootPath string
	markdown goldmark.Markdown
}

func NewStore(rootPath string) *Store {
	return &Store{
		rootPath: rootPath,
		markdown: goldmark.New(),
	}
}

// Parse parses the localized variant of a post.
//
// If not found, `store.NotFoundError` is returned.
func (s *Store) Parse(path string, locale string) (*PostData, error) {
	contents, err := s.readMarkdown(path, locale)
	if err != nil {
		return nil, err
	}

	html, err := s.convertToHtml(contents)
	if err != nil {
		return nil, fmt.Errorf("could not convert file: %w", err)
	}

	var result = &PostData{HtmlContent: template.HTML(html)}
	result.populateMetadata(contents)
	return result, nil
}

func (s *Store) convertToHtml(contents []byte) (html string, err error) {
	buffer := bytes.Buffer{}
	err = s.markdown.Convert(contents, &buffer)
	if err == nil {
		html = buffer.String()
	}
	return
}

func (s *Store) readMarkdown(path string, locale string) ([]byte, error) {
	root, err := os.OpenRoot(s.rootPath)
	if err != nil {
		return nil, fmt.Errorf("could not open root: %w", err)
	}
	defer func() { _ = root.Close() }()

	file, err := localization.FindFile(root, path)
	if err != nil {
		return nil, newFileNotFoundError(path, err)
	}

	resolvedPath, isFound := file.UseLocale(locale)
	if !isFound {
		return nil, newLocaleNotFoundError(path, locale)
	}

	return root.ReadFile(resolvedPath)
}

type PostData struct {
	HtmlContent template.HTML
}

func (s *PostData) populateMetadata(contents []byte) {
	// todo: figure out and populate metadata
}
