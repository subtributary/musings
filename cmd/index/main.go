package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/subtributary/musings/internal/posts"
)

const (
	ContentPath = "./content/"
	DataPath    = "./data/"
)

func main() {
	var locale, target string
	flag.StringVar(&locale, "locale", "", "")
	flag.Usage = printUsage
	flag.Parse()
	if flag.NArg() > 1 {
		printUsage()
		os.Exit(1)
	}
	target = flag.Arg(0)

	index, err := posts.LoadIndex(DataPath, locale)
	if err != nil {
		log.Printf("Error loading index: %v", err)
		if index, err = posts.NewIndex(); err != nil {
			log.Fatalf("Error creating new index: %v", err)
		}
		log.Println("Created new index. Continuing.")
	}

	contentPath := filepath.Join(ContentPath, locale)
	root, err := os.OpenRoot(contentPath)
	if err != nil {
		log.Fatalf("Error opening root: %v", err)
	}
	defer func() { _ = root.Close() }()

	if target == "" {
		indexDirectory(index, root.FS())
	} else if rel, err := filepath.Rel(contentPath, target); err != nil {
		log.Fatal("Target path is not within the content directory.")
	} else if err := indexFile(index, root.FS(), rel); err != nil {
		log.Fatalf("Error indexing file: %v", err)
	}

	if err := index.SaveIndex(DataPath, locale); err != nil {
		log.Fatalf("Error saving index: %v", err)
	}
}

func printUsage() {
	program := filepath.Base(os.Args[0])
	fmt.Println("Musings post indexer.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s [options]\n", program)
	fmt.Printf("  %s [options] <file>\n", program)
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --locale <tag>  Set the locale for the index. [default: none]")
}

func indexDirectory(index posts.Index, dir fs.FS) {
	err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if info, err := d.Info(); err != nil {
			return fmt.Errorf("file info: %v", err)
		} else if !info.Mode().IsRegular() {
			return nil
		}
		return indexFile(index, dir, path)
	})
	if err != nil {
		log.Fatalf("Error indexing directory: %v", err)
	}
}

func indexFile(index posts.Index, dir fs.FS, path string) error {
	parser := posts.NewParser()
	if post, err := parser.ParseFile(dir, path); err != nil {
		return fmt.Errorf("parse post: %w", err)
	} else if err = index.Upsert(path, post); err != nil {
		return fmt.Errorf("upsert: %w", err)
	}
	return nil
}
