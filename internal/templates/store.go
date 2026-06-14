package templates

type Store interface {
	Lookup(name string) (Template, error)
}

func NewStore(rootPath string, funcs Funcs, live bool) (Store, error) {
	if live {
		return NewLiveStore(rootPath, funcs), nil
	}

	store := NewCachedStore(funcs)
	err := store.Load(rootPath)
	return store, err
}
