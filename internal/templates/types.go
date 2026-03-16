package templates

import "html/template"

type Store struct {
	// template = map[locale]
	templates map[string]*template.Template
}
