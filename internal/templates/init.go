package templates

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// NewStore creates a template store from the ".gohtml" templates in the path.
// Files should be named like "{name}.gohtml" or "{name}.{locale}.gohtml".
func NewStore(templatesPath string) (*Store, error) {
	store := &Store{
		templates: make(localizedTemplate),
	}

	partialsPath := filepath.Join(templatesPath, PartialsPath)
	if err := store.parseFiles("partials", partialsPath); err != nil {
		return nil, fmt.Errorf("error parsing partial templates: %w", err)
	}

	pagesPath := filepath.Join(templatesPath, PagesPath)
	if err := store.parseFiles("pages", pagesPath); err != nil {
		return nil, fmt.Errorf("error parsing page templates: %w", err)
	}

	return store, nil
}

func (s *Store) parseFiles(category string, path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read templates: %w", err)
	}

	for _, file := range files {
		filename := file.Name()
		if file.IsDir() || filepath.Ext(filename) != ".gohtml" {
			continue
		}

		contents, err := os.ReadFile(filepath.Join(path, filename))
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", filename, err)
		}

		name, locale := ExtractNameAndLocale(filename)
		name = strings.Join([]string{category, name}, "/")

		if s.templates[locale] == nil {
			s.templates[locale] = template.New("")
		}
		if _, err := s.templates[locale].New(name).Parse(string(contents)); err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}
	}

	return nil
}
