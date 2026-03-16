package templates

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// Architectural notes:
// todo for v0.2: right now, the pages determine the supported locales, but this needs to be defined by the config.
// This would also allow "index.gohtml" to reference "partials/header.ko.gohtml", which currently cannot happen.
// I don't want to fix this in v0.1 since I'm not focused on localization here.
// (Let's not put the cart before the horse.)

func New(templatesPath string) (*Store, error) {
	store := &Store{
		templates: make(map[string]*template.Template),
	}

	pagesPath := templatesPath
	err := store.loadRootPages(pagesPath)
	if err != nil {
		return nil, fmt.Errorf("error loading page templates: %w", err)
	}

	partialsPath := filepath.Join(templatesPath, "partials")
	err = store.loadPartials(partialsPath)
	if err != nil {
		return nil, fmt.Errorf("error loading partials: %w", err)
	}

	return store, nil
}

func (s *Store) loadRootPages(path string) error {
	files, err := loadFiles(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if s.templates[file.locale] == nil {
			s.templates[file.locale] = template.New("")
		}
		if _, err := s.templates[file.locale].New(file.name).Parse(file.content); err != nil {
			return fmt.Errorf("error parsing page template %s: %w", file.name, err)
		}
	}

	return nil
}

func (s *Store) loadPartials(path string) error {
	partials, err := loadFiles(path)
	if err != nil {
		return err
	}

	// Map name->locale->content so we can easily address localized variants.
	partialsMap := make(map[string]map[string]string)
	for _, file := range partials {
		if partialsMap[file.name] == nil {
			partialsMap[file.name] = make(map[string]string)
		}
		partialsMap[file.name][file.locale] = file.content
	}

	for locale, tmpl := range s.templates {
		for name, variants := range partialsMap {
			content := findBestVariant(variants, locale)
			if _, err := tmpl.New("partials/" + name).Parse(content); err != nil {
				return fmt.Errorf("error parsing partial template %s: %w", name, err)
			}
		}
	}

	return nil
}

func findBestVariant(variants map[string]string, locale string) string {
	// todo for v0.2: fallback to less specific locales
	priority := []string{locale, ""}

	for _, p := range priority {
		if content, ok := variants[p]; ok {
			return content
		}
	}

	// todo for v0.2: allow partials to not exist for some locales
	panic(fmt.Sprintf("no template variant found for locale %q", locale))
}

type templateFile struct {
	locale  string
	name    string
	content string
}

// loadFiles reads all files in a directory that are named per convention.
func loadFiles(path string) ([]templateFile, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error scanning templates: %w", err)
	}

	result := make([]templateFile, 0, len(files))

	for _, file := range files {
		filename := file.Name()
		if file.IsDir() || filepath.Ext(filename) != ".gohtml" {
			continue
		}

		entry := templateFile{}

		lastDot := len(filename) - len(".gohtml")
		penultimateDot := strings.LastIndex(filename[:lastDot], ".")
		if penultimateDot != -1 {
			entry.name = filename[:penultimateDot]
			entry.locale = strings.ToLower(filename[penultimateDot+1 : lastDot])
		} else {
			entry.name = filename[:lastDot]
			entry.locale = ""
		}

		content, err := os.ReadFile(filepath.Join(path, filename))
		if err != nil {
			return nil, fmt.Errorf("error reading template: %w", err)
		}
		entry.content = string(content)

		result = append(result, entry)
	}

	return result, nil
}
