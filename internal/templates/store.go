package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

type Store struct {
	templates *template.Template
}

func NewStore(templateFS fs.FS, contentFS fs.FS, staticFS fs.FS) (Store, error) {
	templates := template.New("")
	Funcs{ContentFS: contentFS, StaticFS: staticFS}.ApplyTo(templates)

	err := loadFS(templates, templateFS)
	if err != nil {
		return Store{}, fmt.Errorf("load templates: %w", err)
	}

	return Store{
		templates: templates,
	}, nil
}

func loadFS(templates *template.Template, templateFS fs.FS) error {
	return fs.WalkDir(templateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() || filepath.Ext(d.Name()) != ".gohtml" {
			return nil
		}

		name := strings.TrimSuffix(path, ".gohtml")

		content, err := fs.ReadFile(templateFS, path)
		if err != nil {
			return fmt.Errorf("read file %s: %w", path, err)
		}

		_, err = templates.New(name).Parse(string(content))
		if err != nil {
			return fmt.Errorf("parse file %s: %w", path, err)
		}

		return nil
	})
}

func (s Store) Lookup(name string) (*template.Template, error) {
	tmpl := s.templates.Lookup(name)
	if tmpl == nil {
		return nil, fmt.Errorf("template not found: %s", name)
	}
	return tmpl, nil
}
