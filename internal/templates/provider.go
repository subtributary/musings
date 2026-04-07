package templates

import (
	"fmt"
)

type Provider interface {
	Get() (*Store, error)
}

type CachedProvider struct {
	store Store
}

func NewCachedProvider(path string) (provider CachedProvider, err error) {
	provider.store = newTemplatesStore()
	err = provider.store.Load(path)
	if err != nil {
		err = fmt.Errorf("load cached templates: %w", err)
	}
}

func (p *CachedProvider) Get() (*Store, error) {
	return p.store, nil
}

// LiveProvider rebuilds the store every time that Get is called.
type LiveProvider struct {
	path string
}

func NewLiveProvider(path string) *LiveProvider {
	return &LiveProvider{path}
}

func (p *LiveProvider) Get() (*Store, error) {
	store := newTemplatesStore()
	if err := store.Load(p.path); err != nil {
		return nil, fmt.Errorf("load lives templates: %w", err)
	}
	return store, nil
}
