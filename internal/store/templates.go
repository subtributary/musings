package store

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/subtributary/musings/internal/localization"
)

type TemplateProvider interface {
	Get() (*TemplatesStore, error)
}

type CachedTemplateProvider struct {
	store *TemplatesStore
}

func NewCachedTemplateProvider(path string) (*CachedTemplateProvider, error) {
	store := newTemplatesStore()
	if err := store.Load(path); err != nil {
		return nil, fmt.Errorf("load cached templates: %w", err)
	}
	return &CachedTemplateProvider{store: store}, nil
}

func (p *CachedTemplateProvider) Get() (*TemplatesStore, error) {
	return p.store, nil
}

// LiveTemplateProvider rebuilds the store every time that Get is called.
type LiveTemplateProvider struct {
	path string
}

func NewLiveTemplateProvider(path string) *LiveTemplateProvider {
	return &LiveTemplateProvider{path}
}

func (p *LiveTemplateProvider) Get() (*TemplatesStore, error) {
	store := newTemplatesStore()
	if err := store.Load(p.path); err != nil {
		return nil, fmt.Errorf("load lives templates: %w", err)
	}
	return store, nil
}

type TemplatesStore struct {
	// templates is a map of locales to templates
	templates map[string]*template.Template
}

func newTemplatesStore() *TemplatesStore {
	return &TemplatesStore{
		templates: make(map[string]*template.Template),
	}
}

// Lookup returns the localized variant of the named template.
//
// If not found, nil is returned.
func (s *TemplatesStore) Lookup(name string, locale string) (tmpl *template.Template) {
	_, _ = localization.FindMatchCb(s.templates, locale, func(t *template.Template) bool {
		tmpl = t.Lookup(name)
		return tmpl != nil
	})
	return
}

// Load parses all templates in a conventional templates folder into memory.
func (s *TemplatesStore) Load(path string) error {
	// todo: pass this in from somewhere
	locales := []string{"en", ""}

	root, err := os.OpenRoot(path)
	if err != nil {
		return fmt.Errorf("open root: %s", err)
	}
	defer func() { _ = root.Close() }()

	if err := s.load(root, ".", locales); err != nil {
		return fmt.Errorf("load page templates: %w", err)
	}
	if err := s.load(root, "partials", locales); err != nil {
		return fmt.Errorf("load partial templates: %w", err)
	}
	return nil
}

func (s *TemplatesStore) load(root *os.Root, path string, locales []string) error {
	files, err := localization.ScanFiles(root, path)
	if err != nil {
		return fmt.Errorf("load templates: %w", err)
	}

	prefix := path + "/"
	if prefix == "./" {
		prefix = ""
	}

	for _, locale := range locales {
		if s.templates[locale] == nil {
			s.templates[locale] = template.New("")
		}
		for name, file := range files {
			if filepath.Ext(name) != ".gohtml" {
				continue
			}

			resolvedPath, isFound := file.UseLocale(locale)
			if !isFound {
				return newLocaleNotFoundError(path, locale)
			}

			contents, err := root.ReadFile(resolvedPath)
			if err != nil {
				return fmt.Errorf("read file %q: %w", resolvedPath, err)
			}

			name = prefix + strings.TrimSuffix(name, ".gohtml")
			if _, err := s.templates[locale].New(name).Parse(string(contents)); err != nil {
				return fmt.Errorf("parse %q: %w", resolvedPath, err)
			}
		}
	}

	return nil
}
