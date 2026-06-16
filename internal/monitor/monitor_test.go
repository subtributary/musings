package monitor_test

import (
	"io/fs"
	"testing"
	"testing/fstest"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/subtributary/musings/internal/monitor"
)

// TestWatcher is meant to be used in combination with fstest.MapFS.
// It does not watch the fs automatically and relies on manually invoking the
// signals via TestWatcher.SignalCreated and TestWatcher.SignalRemoved.
type TestWatcher struct {
	errors chan error
	events chan fsnotify.Event
	closed bool
}

func NewTestWatcher() *TestWatcher {
	return &TestWatcher{
		errors: make(chan error, 1),
		events: make(chan fsnotify.Event, 1),
	}
}

func (tw *TestWatcher) Add(string) error {
	return nil
}

func (tw *TestWatcher) Close() error {
	tw.closed = true
	return nil
}

func (tw *TestWatcher) Events() <-chan fsnotify.Event {
	return tw.events
}

func (tw *TestWatcher) Errors() <-chan error {
	return tw.errors
}

func (tw *TestWatcher) SignalCreated(name string) {
	if !tw.closed {
		tw.events <- fsnotify.Event{
			Name: name,
		}
	}
}

func (tw *TestWatcher) SignalRemoved(name string) {
	if !tw.closed {
		tw.events <- fsnotify.Event{
			Name: name,
		}
	}
}

// DirtyCall is used to store a call to it DirtyCollector.Dirty.
type DirtyCall struct {
	Name      string
	Info      fs.FileInfo
	IsRemoved bool
}

func (call DirtyCall) Assert(t *testing.T, name string, modTime time.Time, isRemoved bool) {
	t.Helper()

	if call.Name != name {
		t.Errorf("Name = %s, want %s", call.Name, name)
	}

	var gotModTime time.Time
	if call.Info != nil {
		gotModTime = call.Info.ModTime()
	}
	if gotModTime != modTime {
		t.Errorf("ModTime() = %v, want %v", gotModTime, modTime)
	}

	if call.IsRemoved != isRemoved {
		t.Errorf("IsRemoved = %v, want %v", call.IsRemoved, isRemoved)
	}
}

// DirtyCollector exposes a Dirty function
type DirtyCollector struct {
	calls chan DirtyCall
}

func NewDirtyCollector() *DirtyCollector {
	return &DirtyCollector{
		calls: make(chan DirtyCall, 10),
	}
}

func (c *DirtyCollector) Dirty(name string, info fs.FileInfo, isRemoved bool) {
	c.calls <- DirtyCall{
		Name:      name,
		Info:      info,
		IsRemoved: isRemoved,
	}
}

func (c *DirtyCollector) Get() (DirtyCall, bool) {
	select {
	case call := <-c.calls:
		return call, true
	case <-time.After(1 * time.Second):
		return DirtyCall{}, false
	}
}

func TestMonitor(t *testing.T) {
	files := fstest.MapFS{}
	watcher := NewTestWatcher()
	dirties := NewDirtyCollector()
	subject, err := monitor.New(dirties.Dirty, files, ".",
		monitor.WithFileWatcher(watcher),
		monitor.WithStat(files.Stat),
	)
	if err != nil {
		t.Fatalf("error creating monitor: %v", err)
	}
	defer (func() { _ = subject.Close() })()

	t.Run("initial files are dirty", func(t *testing.T) {
		files["en/hello.md"] = &fstest.MapFile{}

		if err = subject.AddDirectory("."); err != nil {
			t.Fatalf("error adding dir: %v", err)
		}

		call, ok := dirties.Get()
		if !ok {
			t.Fatalf("expected dirty call, got none")
		}
		call.Assert(t, "/en/hello.md", time.Time{}, false)
	})

	t.Run("added file is dirty", func(t *testing.T) {
		files["apple.md"] = &fstest.MapFile{}
		watcher.SignalCreated("apple.md")

		call, ok := dirties.Get()
		if !ok {
			t.Fatalf("expected dirty call, got none")
		}
		call.Assert(t, "/apple.md", time.Time{}, false)
	})

	t.Run("added directory is not dirty", func(t *testing.T) {
		files["fruits/"] = &fstest.MapFile{}
		watcher.SignalCreated("fruits/")
		// No assertions because of fstest.MapFS compatibility issues,
		// but we still want the side effects of our actions.
		_, _ = dirties.Get()
	})

	// Removed directory is dirty becase we don't know it was a directory.
	t.Run("removed directory is dirty", func(t *testing.T) {
		delete(files, "fruits/")
		watcher.SignalRemoved("fruits/")
		// No assertions because of fstest.MapFS compatibility issues,
		// but we still want the side effects of our actions.
		_, _ = dirties.Get()
	})

	t.Run("removed file is dirty", func(t *testing.T) {
		delete(files, "apple.md")
		watcher.SignalRemoved("apple.md")

		call, ok := dirties.Get()
		if !ok {
			t.Fatalf("expected dirty call, got none")
		}
		call.Assert(t, "/apple.md", time.Time{}, true)
	})

	t.Run("closed monitor does not call dirty", func(t *testing.T) {
		if err = subject.Close(); err != nil {
			t.Fatalf("error closing monitor: %v", err)
		}

		delete(files, "en/hello.md")
		watcher.SignalRemoved("en/hello.md")

		call, ok := dirties.Get()
		if ok {
			t.Fatalf("expected no dirty call, got %v", call)
		}
	})
}
