package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
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

func NewCachedStore(rootPath string, tags []language.Tag) (s CachedStore, err error) {
	s.templates = make(map[language.Tag]*template.Template)
	for _, tag := range tags {
		s.templates[tag] = template.New("")
	}

	root, err := os.OpenRoot(rootPath)
	if err != nil {
		err = fmt.Errorf("open root: %s", err)
		return
	}
	defer func() { _ = root.Close() }()

	if err = s.loadTemplatesFromPath(root, ".", tags, ""); err != nil {
		err = fmt.Errorf("load page templates: %w", err)
	} else if err = s.loadTemplatesFromPath(root, "partials", tags, "partials/"); err != nil {
		err = fmt.Errorf("load partial templates: %w", err)
	}

	return
}

func (s CachedStore) Execute(w http.ResponseWriter, name string, tag language.Tag, data any) error {
	if tmpl := s.lookupTemplate(name, tag); tmpl == nil {
		return fmt.Errorf("template for locale %v not found", tag)
	} else {
		return tmpl.Execute(w, data)
	}
}

func (s CachedStore) lookupTemplate(name string, tag language.Tag) *template.Template {
	tmpl := s.templates[tag]
	if tmpl != nil {
		tmpl = tmpl.Lookup(name)
	}
	return tmpl
}

func (s CachedStore) loadTemplatesFromPath(root *os.Root, path string, tags []language.Tag, prefix string) error {
	dir, err := root.OpenRoot(path)
	if err != nil {
		return fmt.Errorf("open root: %w", err)
	}
	defer func() { _ = dir.Close() }()

	scanResult, err := localization.Scan(dir.FS())
	if err != nil {
		return fmt.Errorf("scan files: %w", err)
	}
	groupedFiles := scanResult.GroupByTag(tags)

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

func NewLiveStore(rootPath string, tags []language.Tag) LiveStore {
	return LiveStore{
		rootPath: rootPath,
		tags:     tags,
	}
}

func (s LiveStore) Execute(w http.ResponseWriter, name string, tag language.Tag, data any) error {
	tmpl, err := NewCachedStore(s.rootPath, s.tags)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, name, tag, data)
}
