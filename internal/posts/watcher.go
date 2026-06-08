package posts

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	root  string
	dirty func()
}

func NewWatcher(root string, dirty func()) Watcher {
	return Watcher{
		root:  root,
		dirty: dirty,
	}
}

func (w *Watcher) Start(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := addDirectories(watcher, w.root); err != nil {
		_ = watcher.Close()
		return err
	}

	go w.runWatcher(ctx, watcher)

	return nil
}

func (w *Watcher) runWatcher(ctx context.Context, watcher *fsnotify.Watcher) {
	defer (func() { _ = watcher.Close() })()

	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Create) {
				w.watchNewDirectory(watcher, event.Name)
			}

			w.dirty()

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Printf("posts watcher: %v", err)
		}
	}
}

func (w *Watcher) watchNewDirectory(watcher *fsnotify.Watcher, path string) {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return
	}

	if err := addDirectories(watcher, path); err != nil {
		log.Printf("watch new directory %q: %v", path, err)
	}
}

func addDirectories(watcher *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return watcher.Add(path)
		}

		return nil
	})
}
