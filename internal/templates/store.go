package templates

import (
	"fmt"
	"html/template"
	"os"

	"github.com/subtributary/musings/internal/localization"
)

type Store struct {
	templates map[string]*template.Template
}

func newStore() *Store {
	return &Store{
		templates: make(map[string]*template.Template),
	}
}

func (s *Store) Lookup(name string, locale string) (tmpl *template.Template) {
	_, _ = localization.FindMatchCb(s.templates, locale, func(t *template.Template) bool {
		tmpl = t.Lookup(name)
		return tmpl != nil
	})
	return
}

func (s *Store) Load(path string) error {
	// todo: pass this in from somewhere
	locales := []string{"en", ""}

	root, err := os.OpenRoot(path)
	if err != nil {
		return fmt.Errorf("open root: %s", err)
	}
	defer func() { _ = root.Close() }()

	if err := s.load(root, ".", locales, ""); err != nil {
		return fmt.Errorf("load page templates: %w", err)
	}
	if err := s.load(root, "partials", locales, "partials/"); err != nil {
		return fmt.Errorf("load partial templates: %w", err)
	}
	return nil
}

func (s *Store) load(root *os.Root, path string, locales []string, prefix string) error {
	files, err := localization.ScanFiles(root, path, ".gohtml")
	if err != nil {
		return fmt.Errorf("scan files: %s", err)
	}

	for _, locale := range locales {
		if s.templates[locale] == nil {
			s.templates[locale] = template.New("")
		}
		for name, file := range files {
			filePath, err := file.UseLocale(locale)
			if err != nil {
				return err
			}
			if err := s.addTemplate(root, locale, prefix+name, filePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Store) addTemplate(root *os.Root, locale, name, filePath string) error {
	contents, err := root.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read %q: %w", filePath, err)
	}
	if _, err := s.templates[locale].New(name).Parse(string(contents)); err != nil {
		return fmt.Errorf("could not parse %q: %w", locale, err)
	}
	return nil
}
