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

	if cfg.TargetFile != "" && len(cfg.TargetLocales) == 1 {
		if err := indexFile(cfg.TargetLocales[0], cfg.TargetFile); err != nil {
			log.Fatalf("Error indexing file: %v", err)
		}
	} else if cfg.TargetFile != "" && len(cfg.TargetLocales) > 1 {
		for _, locale := range cfg.TargetLocales {
			if err := indexFile(locale, cfg.TargetFile); err != nil {
				log.Printf("Error indexing file for locale %q: %v", locale, err)
			}
		}
	} else if cfg.TargetFile == "" {
		for _, locale := range cfg.TargetLocales {
			if err := indexDir(locale); err != nil {
				log.Printf("Error indexing locale %q: %v", locale, err)
			}
		}
	}
}

func indexDir(locale language.Tag) error {
	index, err := posts.NewIndex()
	if err != nil {
		return fmt.Errorf("new index: %w", err)
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
			return fmt.Errorf("parse %q: %w", path, err)
		}

		if err := index.Upsert(path, post); err != nil {
			return fmt.Errorf("upsert %q: %w", path, err)
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
		return fmt.Errorf("load index: %w", err)
	}

	contentRoot := os.DirFS(config.LocalizedContentPath(locale))
	post, err := posts.NewParser().ParseFile(contentRoot, path)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	if err := index.Upsert(path, post); err != nil {
		return fmt.Errorf("upsert: %w", err)
	}

	if err := index.SaveIndex(config.DataPath, locale); err != nil {
		return fmt.Errorf("save index: %w", err)
	}
	return nil
}
