package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type CachedStore struct {
	templates *template.Template
}

func NewCachedStore(funcs Funcs) *CachedStore {
	store := &CachedStore{
		templates: template.New(""),
	}

	funcs.ApplyTo(store.templates)

	return store
}

func (s *CachedStore) Load(rootPath string) error {
	root, err := os.OpenRoot(rootPath)
	if err != nil {
		return fmt.Errorf("open root: %w", err)
	}
	defer func() { _ = root.Close() }()

	if err = s.LoadFS(root.FS()); err != nil {
		return fmt.Errorf("load templates: %w", err)
	}

	return nil
}

func (s *CachedStore) LoadFS(dir fs.FS) error {
	return fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() || filepath.Ext(d.Name()) != ".gohtml" {
			return nil
		}

		// Template name is just its path minus its extension
		name := strings.TrimSuffix(path, ".gohtml")

		contents, err := fs.ReadFile(dir, path)
		if err != nil {
			return fmt.Errorf("read file %q: %w", path, err)
		}

		_, err = s.templates.New(name).Parse(string(contents))
		if err != nil {
			return fmt.Errorf("parse %q: %w", path, err)
		}
		return nil
	})
}

func (s *CachedStore) Lookup(name string) (tmpl *template.Template, err error) {
	tmpl = s.templates.Lookup(name)
	if tmpl == nil {
		err = fmt.Errorf("template not found: %s", name)
	}
	return
}
