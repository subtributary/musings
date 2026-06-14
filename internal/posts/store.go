package posts

import (
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"log"
	"os"
	"path"
	"time"
)

// Store ties all the functionality of the package together.
type Store struct {
	index    *Index
	modTimes map[string]time.Time
	parser   Parser
	root     *os.Root
	watcher  *Watcher
}

func OpenStore(contentRoot string) (s Store, err error) {
	s.index = NewIndex()
	s.modTimes = make(map[string]time.Time)
	s.parser = NewParser()

	root, err := os.OpenRoot(contentRoot)
	if err != nil {
		return Store{}, fmt.Errorf("open content root: %v", err)
	}
	s.root = root

	// Start the watcher last because it will start building the index.
	s.watcher = NewWatcher(contentRoot, s.dirty)
	if err = s.watcher.Start(); err != nil {
		err = fmt.Errorf("start watching content root: %v", err)
		return Store{}, errors.Join(err, s.Close())
	}

	return s, nil
}

func (s *Store) Close() (err error) {
	if s.root != nil {
		err = errors.Join(err, s.root.Close())
		s.root = nil
	}
	if s.watcher != nil {
		err = errors.Join(err, s.watcher.Close())
		s.watcher = nil
	}
	return err
}

func (s *Store) Post(path string) (ParsedPost, error) {
	post, err := s.parser.ParseFile(s.root.FS(), path)
	if err != nil {
		return post, err
	}

	modified, _ := s.modTimes[path]
	post.Modified = modified
	return post, nil
}

func (s *Store) Search(query string) iter.Seq[IndexedPost] {
	return s.index.Search(query)
}

func (s *Store) dirty(name string) {
	info, err := fs.Stat(s.root.FS(), name)

	if err != nil {
		s.index.Remove(name)
		delete(s.modTimes, name)
		return
	}

	if path.Ext(info.Name()) == ".md" {
		post, err := s.parser.ParseFile(s.root.FS(), name)
		if err != nil {
			log.Printf("Error parsing post %q: %v", name, err)
			return
		}
		s.index.Upsert(name, post)
	}

	s.modTimes[name] = info.ModTime()
}
