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

type Watcher struct {
	dirty   func(name string)
	rootDir string
	watcher *fsnotify.Watcher
}

// NewWatcher creates a new watcher for a directory and its subdirectories.
// The dirty function is called at Watcher.Start and when a file is modified.
func NewWatcher(rootDir string, dirty func(name string)) *Watcher {
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
	info, err := os.Stat(event.Name)
	switch {
	case err != nil:
		w.dirty(event.Name)
	case !info.IsDir():
		w.dirty(event.Name)
	case event.Has(fsnotify.Create):
		_ = w.addDirectory(event.Name)
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

		path, err = filepath.Rel(w.rootDir, path)
		if err != nil {
			return fmt.Errorf("relative path: %w", err)
		}
		w.dirty(path)

		return nil
	})
}
