package posts

import (
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type AutoIndex struct {
	parser   Parser
	rootPath string
	watcher  *fsnotify.Watcher
	wrapped  *Index
}

func NewAutoIndex(rootPath string) (AutoIndex, error) {
	// URL versioning is not needed for searching, so use an identity function.
	versionURL := func(_, target string) string { return target }

	idx := AutoIndex{
		parser:   NewParser(versionURL),
		rootPath: rootPath,
		wrapped:  NewIndex(),
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return idx, fmt.Errorf("create watcher: %w", err)
	}
	idx.watcher = watcher

	err = idx.watchDir(rootPath)
	if err != nil {
		err = errors.Join(err, idx.Close())
		return idx, err
	}

	return idx, nil
}

func (idx AutoIndex) Close() error {
	return idx.watcher.Close()
}

// List lists all posts in descending order of publication date.
// Posts with publication dates in the future are omitted.
func (idx AutoIndex) List() iter.Seq[*IndexedPost] {
	return idx.wrapped.List()
}

// Search returns the posts matching the query,
// sorted by match score with the best match first.
// Posts with publication dates in the future are omitted.
func (idx AutoIndex) Search(query string) iter.Seq[*IndexedPost] {
	return idx.wrapped.Search(query)
}

func (idx AutoIndex) key(name string) string {
	name = strings.ReplaceAll(name, string(os.PathSeparator), "/")
	name, _ = strings.CutPrefix(name, idx.rootPath)
	name, _ = strings.CutPrefix(name, "/")
	name = "/" + name
	return name
}

func (idx AutoIndex) handleEvent(event fsnotify.Event) error {
	info, err := os.Stat(event.Name)

	// If info error, the file or directory was removed
	if err != nil {
		idx.wrapped.Remove(idx.key(event.Name))
		return nil
	}

	// If it is a directory and was just created, we want to watch it.
	if info.IsDir() && event.Has(fsnotify.Create) {
		err = idx.watchDir(event.Name)
		if err != nil {
			return fmt.Errorf("watch directory: %w", err)
		}
	}

	// If it's a post, update it in the index.
	if !info.IsDir() && path.Ext(event.Name) == ".md" {
		post, err := idx.parser.ParseFile(os.DirFS("."), event.Name)
		if err != nil {
			return fmt.Errorf("parse post: %w", err)
		}
		idx.wrapped.Upsert(idx.key(event.Name), post)
	}

	return nil
}

func (idx AutoIndex) listen() {
	for {
		select {
		case event, ok := <-idx.watcher.Events:
			if !ok {
				return
			}
			if err := idx.handleEvent(event); err != nil {
				log.Printf("Error: %v", err)
			}
		case err, ok := <-idx.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Error: watcher: %v", err)
		}
	}
}

func (idx AutoIndex) watchDir(name string) error {
	return filepath.WalkDir(name, func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			err = idx.watcher.Add(name)
			if err != nil {
				return fmt.Errorf("watch %s: %w", name, err)
			}
			return nil
		}

		if path.Ext(d.Name()) != ".md" {
			return nil
		}

		post, err := idx.parser.ParseFile(os.DirFS("."), name)
		if err != nil {
			return fmt.Errorf("parse post: %w", err)
		}
		idx.wrapped.Upsert(idx.key(name), post)

		return nil
	})
}
