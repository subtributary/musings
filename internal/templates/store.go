package templates

type Store interface {
	Lookup(name string) (Template, error)
}

func NewStore(rootPath string, live bool) (Store, error) {
	if live {
		return NewLiveStore(rootPath), nil
	}

	store := NewCachedStore()
	err := store.Load(rootPath)
	return store, err
}
