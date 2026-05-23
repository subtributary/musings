package templates

import (
	"fmt"
	"os"
)

type LiveStore struct {
	rootPath string
}

func NewLiveStore(rootPath string) LiveStore {
	return LiveStore{rootPath: rootPath}
}

func (s LiveStore) Lookup(name string) (Template, error) {
	root, err := os.OpenRoot(s.rootPath)
	if err != nil {
		return Template{}, fmt.Errorf("open root: %w", err)
	}
	defer func() { _ = root.Close() }()

	store := NewCachedStore()
	if err = store.LoadFS(root.FS()); err != nil {
		return Template{}, fmt.Errorf("load store: %w", err)
	}

	return store.Lookup(name)
}
