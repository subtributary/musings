package templates

import "fmt"

type TemplateProvider interface {
	Get() (*Store, error)
}

type CachedTemplateProvider struct {
	store *Store
}

func NewCachedTemplateProvider(path string) (*CachedTemplateProvider, error) {
	store := newStore()
	if err := store.Load(path); err != nil {
		return nil, fmt.Errorf("load cached templates: %w", err)
	}
	return &CachedTemplateProvider{store: store}, nil
}

func (p *CachedTemplateProvider) Get() (*Store, error) {
	return p.store, nil
}

// LiveTemplateProvider rebuilds the store every time that Get is called.
type LiveTemplateProvider struct {
	path string
}

func NewLiveTemplateProvider(path string) *LiveTemplateProvider {
	return &LiveTemplateProvider{path}
}

func (p *LiveTemplateProvider) Get() (*Store, error) {
	store := newStore()
	if err := store.Load(p.path); err != nil {
		return nil, fmt.Errorf("load lives templates: %w", err)
	}
	return store, nil
}
