package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"slices"
	"strings"
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

func OpenContent(contentRoot *os.Root, locales []localization.Locale) (*Content, error) {
	c := &Content{
		root: contentRoot,
	}

	if len(locales) == 0 {
		locales = []localization.Locale{localization.UndLocale}
	}
	c.locales = locales

	c.indexes = make(map[string]*posts.Index, len(locales))
	for _, loc := range locales {
		c.indexes[loc.Tag] = posts.NewIndex()
	}

	c.parser = posts.NewParser(c.ModTime)

	m, err := monitor.New(c.dirty, contentRoot.FS(), ContentPath)
	if err != nil {
		err = fmt.Errorf("new monitor: %v", err)
		return nil, errors.Join(err, c.Close())
	}
	c.monitor = m

	// The intitial call blocks.
	err = c.monitor.AddDirectory(".")
	if err != nil {
		err = fmt.Errorf("monitor content directory: %w", err)
		return nil, errors.Join(err, c.Close())
	}

	// Make the indexes thread-safe after initial population.
	for _, idx := range c.indexes {
		if err = idx.MakeConcurrent(); err != nil {
			err = fmt.Errorf("make indexes thread safe: %w", err)
			return nil, errors.Join(err, c.Close())
		}
	}

	// Start listening for updates now that our index is set up.
	go c.monitor.Listen()

	return c, nil
}

func (c *Content) Close() (err error) {
	// Close watcher first because event handling may use root.
	if c.monitor != nil {
		err = c.monitor.Close()
		c.monitor = nil
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

type SearchResults struct {
	Query   string
	Results []posts.IndexedPost
}

func (c *Content) Search(locale localization.Locale, query string) SearchResults {
	index, ok := c.indexes[locale.Tag]
	if !ok {
		log.Printf("Unexpected: missing localized index: %s", locale.Tag)
		return SearchResults{}
	}

	return SearchResults{
		Query:   query,
		Results: slices.Collect(index.Search(query)),
	}
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
	name = strings.TrimSuffix(name, ".md")

	locale, _ := localization.ExtractLocale(c.locales, name)

	index, ok := c.indexes[locale.Tag]
	if !ok {
		log.Printf("Unexpected: localized index is missing: %s", locale.Tag)
		return
	}

	if isRemoved {
		index.Remove(name)
	} else {
		post, err := c.parser.ParseFile(c.root.FS(), name+".md")
		if err != nil {
			log.Printf("Error: parse post %q: %v", name, err)
			return
		}
		index.Upsert(name, post)
	}
}
