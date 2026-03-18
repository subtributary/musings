package files

import (
	"fmt"
	"html/template"
	"log"
	"os"
)

type TemplateStore struct {
	templates map[string]*template.Template
}

func NewTemplateStore() *TemplateStore {
	return &TemplateStore{
		templates: make(map[string]*template.Template),
	}
}

func (s *TemplateStore) Load(path string) error {
	// todo: pass this in from somewhere
	locales := []string{"en", ""}

	root, err := os.OpenRoot(path)
	if err != nil {
		return fmt.Errorf("open root %q: %w", path, err)
	}
	defer func() { _ = root.Close() }()

	pageFiles, err := scan(root, ".", ".gohtml")
	if err != nil {
		return fmt.Errorf("could not scan pages: %w", err)
	}
	partialFiles, err := scan(root, "partials", ".gohtml")
	if err != nil {
		return fmt.Errorf("could not scan partials: %w", err)
	}

	for _, locale := range locales {
		s.templates[locale] = template.New("")
		for name, variants := range pageFiles {
			filePath, err := findBestVariant(variants, locale)
			if err != nil {
				return fmt.Errorf("could not resolve template %s: %w", name, err)
			}
			if err := s.addTemplate(root, locale, name, filePath); err != nil {
				return err
			}
		}
		for name, variants := range partialFiles {
			filePath, err := findBestVariant(variants, locale)
			if err != nil {
				// partials are allowed to be omitted, but we should log this for diagnostics
				log.Printf("could not resolve template %s for locale %s", name, locale)
				continue
			}
			if err := s.addTemplate(root, locale, "partials/"+name, filePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *TemplateStore) Lookup(name string, locale string) *template.Template {
	// todo: get the priorities from elsewhere
	priority := []string{locale, ""}

	for _, p := range priority {
		if collection, ok := s.templates[p]; ok {
			if tmpl := collection.Lookup(name); tmpl != nil {
				return tmpl
			}
		}
	}

	return nil
}

func (s *TemplateStore) addTemplate(root *os.Root, locale string, name string, filePath string) error {
	contents, err := root.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read template %s: %w", filePath, err)
	}
	if _, err := s.templates[locale].New(name).Parse(string(contents)); err != nil {
		return fmt.Errorf("could not load template %s: %w", name, err)
	}
	return nil
}
