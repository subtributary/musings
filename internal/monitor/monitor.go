package monitor

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// DirtyFunc is called when Monitor initializes and when it detects a created,
// removed, or updated file. The `name` argument is the path relative to the
// directory being watched and is prefixed with a "/". The `info` argument is
// only valid if `isRemoved` is false.
type DirtyFunc func(name string, info fs.FileInfo, isRemoved bool)

type Option func(m *Monitor)

func WithFileWatcher(watcher Watcher) Option {
	return func(m *Monitor) {
		m.watcher = watcher
	}
}

type StatFunc func(string) (fs.FileInfo, error)

func WithStat(stat StatFunc) Option {
	return func(m *Monitor) {
		m.stat = stat
	}
}

// Monitor monitors a directory calls a DirtyFunc at initialization and when
// a change is detected.
type Monitor struct {
	dirty    DirtyFunc
	root     fs.FS
	rootPath string
	watcher  Watcher
	stat     func(string) (fs.FileInfo, error)
}

// New creates a new Monitor for a directory.
func New(dirty DirtyFunc, root fs.FS, rootPath string, opts ...Option) (*Monitor, error) {
	m := &Monitor{
		dirty:    dirty,
		root:     root,
		rootPath: path.Clean(rootPath),
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.watcher == nil {
		watcher, err := newFileWatcher()
		if err != nil {
			return nil, err
		}
		m.watcher = watcher
	}

	if m.stat == nil {
		m.stat = os.Stat
	}

	return m, nil
}

// AddDirectory adds a directory and its subdirectories to the monitor. The
// directory must be relative to the root of the directory being watched.
func (m *Monitor) AddDirectory(name string) error {
	return fs.WalkDir(m.root, name, func(relPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			err = m.watcher.Add(path.Join(m.rootPath, relPath))
			if err != nil {
				return fmt.Errorf("watch %s: %w", relPath, err)
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("path info: %w", err)
		}

		m.dirty("/"+relPath, info, false)

		return nil
	})
}

func (m *Monitor) Close() error {
	return m.watcher.Close()
}

func (m *Monitor) Listen() {
	for {
		select {
		case event, ok := <-m.watcher.Events():
			if !ok {
				return
			}
			m.handleEvent(event)
		case err, ok := <-m.watcher.Errors():
			if !ok {
				return
			}
			log.Printf("Error watching: %v\n", err)
		}
	}
}

func (m *Monitor) handleEvent(event fsnotify.Event) {
	info, err := m.stat(event.Name)

	// If info error, the file or directory was removed.
	if err != nil {
		m.signalDirty(event.Name, nil, true)
		return
	}

	// At this point, we know it's an existing file or directory.
	// If it's not a directory, the file was created or updated.
	if !info.IsDir() {
		m.signalDirty(event.Name, info, false)
		return
	}

	// At this point, we know it's an existing directory.
	// If it was just created, we want to watch it.
	if event.Has(fsnotify.Create) {
		if err = m.AddDirectory(event.Name); err != nil {
			log.Printf("Unexpected: cannot watch %q: %v", event.Name, err)
		}
	}
}

func (m *Monitor) signalDirty(name string, info fs.FileInfo, isRemoved bool) {
	name, _ = strings.CutPrefix(name, m.rootPath)
	name, _ = strings.CutPrefix(name, "/")
	name = "/" + name
	m.dirty(name, info, isRemoved)
}
