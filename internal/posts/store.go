package posts

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"

	"github.com/yuin/goldmark"
	"golang.org/x/text/language"
)

type Store struct {
	markdown goldmark.Markdown
	files    files.Store
}

func NewStore(rootPath string, tags []language.Tag) *Store {
	return &Store{
		markdown: goldmark.New(),
		files:    files.NewStore(rootPath, tags),
	}
}

// Parse parses the localized variant of a post.
//
// If not found, `store.NotFoundError` is returned.
func (s *Store) Parse(path string, locale string) (PostData, error) {
	root, err := os.OpenRoot(s.rootPath)
	if err != nil {
		return nil, fmt.Errorf("could not open root: %w", err)
	}
	defer func() { _ = root.Close() }()

	// Read file contents.
	bestTag, _ := language.MatchStrings(s.matcher, locale)
	localizedFS := files.NewLocalizedFS(root.FS(), bestTag)
	contents, err := fs.ReadFile(localizedFS, path)
	if err != nil {
		return nil, fmt.Errorf("could not read file %q: %w", path, err)
	}

	// Convert to HTML.
	buffer := bytes.Buffer{}
	err = s.markdown.Convert(contents, &buffer)
	if err != nil {
		return nil, fmt.Errorf("could not convert file %q: %w", path, err)
	}
	html := buffer.String()

	return PostData{
		HtmlContent: template.HTML(html),
	}, nil
}

type PostData struct {
	HtmlContent template.HTML
}
