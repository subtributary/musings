package templates

import "html/template"

type Store struct {
	templates localizedTemplate
}

// template = map[locale]
type localizedTemplate map[string]*template.Template
