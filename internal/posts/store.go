package posts

import (
	"context"
	"fmt"
	"io/fs"
	"iter"
	"log"
	"os"
	"sync/atomic"
)

type Store struct {
	index       atomic.Pointer[Index]
	contentRoot fs.FS
	wake        chan struct{}
}

// OpenStore opens a store that stays in sync with the file system.
func OpenStore(ctx context.Context, contentRoot string) (*Store, error) {
	s := &Store{
		contentRoot: os.DirFS(contentRoot),
		wake:        make(chan struct{}, 1),
	}

	s.index.Store(NewIndex())

	watcher := NewWatcher(contentRoot, s.MarkDirty)
	if err := watcher.Start(ctx); err != nil {
		return nil, fmt.Errorf("watch dir: %w", err)
	}

	// Start the indexer in the background.
	// Call this after no more errors can happen.
	go s.runIndexer(ctx)

	// Start the initial build.
	s.MarkDirty()

	return s, nil
}

// List lists all posts in order of publication with the most recent one first.
func (s *Store) List() iter.Seq[IndexedPost] {
	snapshot := s.index.Load()
	return snapshot.List()
}

// Lookup returns the post with the given name.
func (s *Store) Lookup(name string) (ParsedPost, error) {
	parser := NewParser()
	return parser.ParseFile(s.contentRoot, name)
}

// MarkDirty signals that the store needs to rebuild its index.
// It is safe to use methods concurrently to the index being rebuilt.
func (s *Store) MarkDirty() {
	select {
	case s.wake <- struct{}{}:
	default:
		// A rebuild is already queued.
	}
}

// Search searches the index per the query string.
func (s *Store) Search(query string) iter.Seq[IndexedPost] {
	snapshot := s.index.Load()
	return snapshot.Search(query)
}

func (s *Store) runIndexer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-s.wake:
			next, err := BuildIndex(ctx, s.contentRoot)
			if err != nil {
				log.Printf("rebuild posts index: %v", err)
				continue
			}
			s.index.Store(next)
		}
	}
}
