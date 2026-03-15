package templates

import (
	"html/template"
	"strings"
)

func (s *Store) Lookup(name string, locale string) *template.Template {
	// Todo: make this not require an exact match
	locale = strings.ToLower(locale)
	collection := s.templates[locale]
	return collection.Lookup(name)
}
