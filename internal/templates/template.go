package templates

import (
	"html/template"
	"io"
	"time"
)

type Template struct {
	wrapped      *template.Template
	lastModified time.Time
}

func (t Template) Execute(w io.Writer, data any) error {
	return t.wrapped.Execute(w, data)
}

func (t Template) LastModified() time.Time {
	return t.lastModified
}
