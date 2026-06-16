package main

import (
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"log"
	"os"
	"path"
	"slices"
	"sync"
	"time"

	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/monitor"
	"github.com/subtributary/musings/internal/posts"
)

type Content struct {
	indexes  map[string]*posts.Index // index := indexes[locale]
	locales  []localization.Locale
	modTimes sync.Map
	monitor  *monitor.Monitor
	root     *os.Root
	parser   posts.Parser
}

func OpenContent(contentRoot string, locales []localization.Locale) (*Content, error) {
	c := &Content{}

	if len(locales) == 0 {
		locales = []localization.Locale{localization.UndLocale}
	}
	c.locales = locales

	c.indexes = make(map[string]*posts.Index, len(locales))
	for _, loc := range locales {
		c.indexes[loc.Tag] = posts.NewIndex()
	}

	c.parser = posts.NewParser(c.ModTime)

	root, err := os.OpenRoot(contentRoot)
	if err != nil {
		return nil, fmt.Errorf("open content root: %v", err)
	}
	c.root = root

	c.monitor, err = monitor.New(c.dirty, root.FS(), contentRoot)
	if err != nil {
		err = fmt.Errorf("new monitor: %v", err)
		return nil, errors.Join(err, c.Close())
	}
	err = c.monitor.AddDirectory(".")
	if err != nil {
		err = fmt.Errorf("monitor content directory: %v", err)
		return nil, errors.Join(err, c.Close())
	}

	return c, nil
}

func (c *Content) Close() (err error) {
	// Close watcher first because event handling may use root.
	if c.monitor != nil {
		err = c.monitor.Close()
		c.monitor = nil
	}

	if c.root != nil {
		err = errors.Join(err, c.root.Close())
		c.root = nil
	}

	return err
}

// GetPost parses a file as a post.
// The name parameter is its path relative to the content root.
func (c *Content) GetPost(name string) (posts.ParsedPost, error) {
	return c.parser.ParseFile(c.root.FS(), name)
}

// ModTime returns the cached modification time of a file.
func (c *Content) ModTime(name string) (time.Time, bool) {
	name = path.Clean("/" + name)

	if t, ok := c.modTimes.Load(name); ok {
		return t.(time.Time), true
	}

	return time.Time{}, false
}

func (c *Content) Search(locale localization.Locale, query string) iter.Seq[posts.IndexedPost] {
	if index, ok := c.indexes[locale.Tag]; ok {
		return index.Search(query)
	}
	log.Printf("Unexpected: missing localized index: %s", locale.Tag)
	return slices.Values([]posts.IndexedPost{})
}

// dirty is called by the monitor when a content file or directory changes.
func (c *Content) dirty(name string, info fs.FileInfo, isRemoved bool) {
	if isRemoved {
		c.modTimes.Delete(name)
	} else {
		c.modTimes.Store(name, info.ModTime())
	}

	// Only continue for files that are likely posts.
	if (!isRemoved && info.IsDir()) || path.Ext(name) != ".md" {
		return
	}

	locale, trailingPath := localization.ExtractLocale(c.locales, name)
	if locale == localization.UndLocale && len(c.locales) != 0 {
		// We are only expecting UndLocale if localization is disabled.
		log.Printf("Unexpected: no index for file locale: %q", name)
		return
	}

	index, ok := c.indexes[locale.Tag]
	if !ok {
		log.Printf("Unexpected: localized index is missing: %s", locale.Tag)
		return
	}

	if isRemoved {
		index.Remove(trailingPath)
	} else {
		post, err := c.parser.ParseFile(c.root.FS(), name)
		if err != nil {
			log.Printf("Error: parse post %q: %v", name, err)
			return
		}
		index.Upsert(trailingPath, post)
	}
}
