package posts

import (
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

// Store ties all the functionality of the package together.
type Store struct {
	index    *Index
	modTimes sync.Map
	parser   Parser
	root     *os.Root
	watcher  *Watcher
}

func OpenStore(contentRoot string) (*Store, error) {
	s := &Store{
		index: NewIndex(),
	}

	s.parser = NewParser(WithModTime(s.modTime))

	root, err := os.OpenRoot(contentRoot)
	if err != nil {
		return nil, fmt.Errorf("open content root: %v", err)
	}
	s.root = root

	// Start the watcher last because it will start building the index.
	s.watcher = NewWatcher(contentRoot, s.dirty)
	if err = s.watcher.Start(); err != nil {
		err = fmt.Errorf("start watching content root: %v", err)
		return nil, errors.Join(err, s.Close())
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

	modified, _ := s.modTimes.Load(path)
	post.Modified = modified.(time.Time)
	return post, nil
}

func (s *Store) Search(query string) iter.Seq[IndexedPost] {
	return s.index.Search(query)
}

func (s *Store) dirty(name string, info fs.FileInfo, isRemoved bool) {
	if isRemoved {
		s.index.Remove(name)
		s.modTimes.Delete(name)
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

	s.modTimes.Store(name, info.ModTime())
}

func (s *Store) modTime(name string) time.Time {
	t, _ := s.modTimes.Load(name)
	return t.(time.Time)
}
