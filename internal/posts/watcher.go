package posts

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type DirtyFunc func(name string, info fs.FileInfo, isRemoved bool)

type Watcher struct {
	dirty   DirtyFunc
	rootDir string
	watcher *fsnotify.Watcher
}

// NewWatcher creates a new watcher for a directory and its subdirectories.
// The dirty function is called at Watcher.Start and when a file is modified.
func NewWatcher(rootDir string, dirty DirtyFunc) *Watcher {
	return &Watcher{
		dirty:   dirty,
		rootDir: rootDir,
	}
}

func (w *Watcher) Start() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}
	w.watcher = watcher

	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				w.handleEvent(event)
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Error watching %q: %v\n", w.rootDir, err)
			}
		}
	}()

	if err = w.addDirectory(w.rootDir); err != nil {
		return errors.Join(err, w.watcher.Close())
	}

	return nil
}

func (w *Watcher) handleEvent(event fsnotify.Event) {
	relPath, err := filepath.Rel(w.rootDir, event.Name)
	if err != nil {
		log.Printf("Unexpected: path is not relative to content root: %v", event.Name)
		return
	}

	info, err := os.Stat(event.Name)

	// If info error, the file or directory was removed.
	if err != nil {
		w.dirty(relPath, info, true)
		return
	}

	// At this point, we know it's an existing file or directory.
	// If it's not a directory, the file was created or updated.
	if !info.IsDir() {
		w.dirty(relPath, info, false)
		return
	}

	// At this point, we know it's an existing directory.
	// If it was just created, we want to watch it.
	if event.Has(fsnotify.Create) {
		if err = w.addDirectory(event.Name); err != nil {
			log.Printf("Unexpected: cannot watch %q: %v", event.Name, err)
		}
	}
}

func (w *Watcher) Close() error {
	if w.watcher != nil {
		return w.watcher.Close()
	}
	return nil
}

func (w *Watcher) addDirectory(name string) error {
	return filepath.WalkDir(name, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return w.watcher.Add(path)
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("path info: %w", err)
		}

		path, err = filepath.Rel(w.rootDir, path)
		w.dirty(path, info, false)

		return nil
	})
}
