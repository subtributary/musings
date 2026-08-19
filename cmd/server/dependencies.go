package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path"

	"github.com/subtributary/musings/internal/app"
)

//go:embed all:templates
var templateFiles embed.FS

//go:embed all:web
var webFiles embed.FS

// Dependencies are shared by multiple routes or require closing.
type Dependencies struct {
	ContentRoot *os.Root
	ViewModels  *ViewModelFactory
	WebFS       fs.FS

	Templates map[string]*template.Template
}

func LoadDependencies(cfg Config) (*Dependencies, error) {
	deps := &Dependencies{}
	var err error

	deps.ContentRoot, err = os.OpenRoot(app.ContentPath)
	if err != nil {
		return nil, fmt.Errorf("open content root: %w", err)
	}

	deps.ViewModels, err = NewModelViewFactory(cfg.Locales)
	if err != nil {
		return nil, fmt.Errorf("create view model factory: %v", err)
	}

	deps.WebFS, err = fs.Sub(webFiles, "web")
	if err != nil {
		return nil, fmt.Errorf("sub embedded web files: %w", err)
	}

	pageFiles, err := fs.ReadDir(templateFiles, "templates/pages")
	if err != nil {
		return nil, fmt.Errorf("list template pages: %w", err)
	}

	deps.Templates = make(map[string]*template.Template)
	for _, f := range pageFiles {
		tmpl, err := template.ParseFS(templateFiles, "templates/*.gohtml")
		if err != nil {
			return nil, fmt.Errorf("load template parts: %w", err)
		}

		_, err = tmpl.ParseFS(templateFiles, path.Join("templates/pages", f.Name()))
		if err != nil {
			return nil, fmt.Errorf("load page template: %w", err)
		}

		deps.Templates[f.Name()] = tmpl
	}

	return deps, nil
}

func (d *Dependencies) Close() error {
	err := d.ContentRoot.Close()
	return err
}
