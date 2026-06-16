package monitor

import "github.com/fsnotify/fsnotify"

type Watcher interface {
	Add(string) error
	Close() error
	Events() <-chan fsnotify.Event
	Errors() <-chan error
}

type fsnotifyWatcher struct {
	*fsnotify.Watcher
}

func (w fsnotifyWatcher) Events() <-chan fsnotify.Event {
	return w.Watcher.Events
}

func (w fsnotifyWatcher) Errors() <-chan error {
	return w.Watcher.Errors
}

func newFileWatcher() (Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return fsnotifyWatcher{watcher}, nil
}
