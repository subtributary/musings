package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/subtributary/musings/internal/localization"
	"golang.org/x/text/language"
)

type Store interface {
	Execute(w http.ResponseWriter, name string, tag language.Tag, data any) error
}

// CachedStore loads and caches templates at startup.
// This is good for efficiency.
type CachedStore struct {
	templates map[language.Tag]*template.Template
}

func NewCachedStore(tags []language.Tag) CachedStore {
	s := CachedStore{}
	for _, tag := range tags {
		s.templates[tag] = template.New("")
	}
	return s
}

func (s CachedStore) Execute(w http.ResponseWriter, name string, tag language.Tag, data any) error {
	tmpl, ok := s.templates[tag]
	if !ok {
		return fmt.Errorf("template for locale %v not found", tag)
	}
	return tmpl.Execute(w, data)
}

func (s CachedStore) Load(rootPath string) error {
	root, err := os.OpenRoot(rootPath)
	if err != nil {
		return fmt.Errorf("open root: %s", err)
	}
	defer func() { _ = root.Close() }()

	if err := s.load(root, ".", ""); err != nil {
		return fmt.Errorf("load page templates: %w", err)
	}
	if err := s.load(root, "partials", "partials/"); err != nil {
		return fmt.Errorf("load partial templates: %w", err)
	}
	return nil
}

func (s CachedStore) load(root *os.Root, path string, prefix string) error {
	dir, err := root.OpenRoot(path)
	if err != nil {
		return fmt.Errorf("open root: %w", err)
	}
	defer func() { _ = dir.Close() }()

	tags := slices.Collect(maps.Keys(s.templates))
	groupedFiles, err := localization.Scan(dir.FS(), tags)
	if err != nil {
		return fmt.Errorf("scan files: %w", err)
	}

	for tag, dirEntries := range groupedFiles {
		localizedFS := localization.NewLocalizedFS(dir.FS(), tag)
		for _, file := range dirEntries {
			name := file.Name()
			if filepath.Ext(name) != ".gohtml" {
				continue
			}

			contents, err := fs.ReadFile(localizedFS, name)
			if err != nil {
				return fmt.Errorf("read file %q: %w", name, err)
			}

			name = prefix + strings.TrimSuffix(name, ".gohtml")
			if _, err := s.templates[tag].New(name).Parse(string(contents)); err != nil {
				return fmt.Errorf("parse %q: %w", name, err)
			}
		}
	}

	return nil
}

// LiveStore loads templates at request time.
// This is good during active development.
type LiveStore struct {
	rootPath string
	tags     []language.Tag
}

func (s LiveStore) Execute(w http.ResponseWriter, name string, tag language.Tag, data any) error {
	temp := NewCachedStore(s.tags)
	err := temp.Load(s.rootPath)
	if err != nil {
		return err
	}
	return temp.Execute(w, name, tag, data)
}
