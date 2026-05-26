package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/subtributary/musings/internal/config"
	"github.com/subtributary/musings/internal/posts"
	"golang.org/x/text/language"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if cfg.Target != "" && len(cfg.Locales) == 1 {
		if err := indexFile(cfg.Locales[0], cfg.Target); err != nil {
			log.Fatalf("Error indexing file: %v", err)
		}
	} else if cfg.Target != "" && len(cfg.Locales) > 1 {
		for _, locale := range cfg.Locales {
			if err := indexFile(locale, cfg.Target); err != nil {
				log.Printf("Error indexing file for locale %q: %v", locale, err)
			}
		}
	} else if cfg.Target == "" {
		for _, locale := range cfg.Locales {
			if err := indexDir(locale); err != nil {
				log.Printf("Error indexing locale %q: %v", locale, err)
			}
		}
	}
}

func indexDir(locale language.Tag) error {
	index, err := posts.NewIndex()
	if err != nil {
		return fmt.Errorf("new index: %v", err)
	}

	contentRoot := os.DirFS(config.LocalizedContentPath(locale))
	err = fs.WalkDir(contentRoot, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() || filepath.Ext(d.Name()) != ".md" {
			return nil
		}

		post, err := posts.NewParser().ParseFile(contentRoot, path)
		if err != nil {
			return fmt.Errorf("parse %q: %v", path, err)
		}

		if err := index.Upsert(path, post); err != nil {
			return fmt.Errorf("upsert %q: %v", path, err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := index.SaveIndex(config.DataPath, locale); err != nil {
		return fmt.Errorf("save index: %v", err)
	}
	return nil
}

func indexFile(locale language.Tag, path string) error {
	index, err := posts.LoadIndex(config.DataPath, locale)
	if err != nil {
		return fmt.Errorf("load index: %v", err)
	}

	contentRoot := os.DirFS(config.LocalizedContentPath(locale))
	post, err := posts.NewParser().ParseFile(contentRoot, path)
	if err != nil {
		return fmt.Errorf("parse file: %v", err)
	}

	if err := index.Upsert(path, post); err != nil {
		return fmt.Errorf("upsert: %v", err)
	}

	if err := index.SaveIndex(config.DataPath, locale); err != nil {
		return fmt.Errorf("save index: %v", err)
	}
	return nil
}
