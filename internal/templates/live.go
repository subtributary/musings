package templates

type LiveStore struct {
	funcs    Funcs
	rootPath string
}

func NewLiveStore(rootPath string, funcs Funcs) LiveStore {
	return LiveStore{funcs: funcs, rootPath: rootPath}
}

func (s LiveStore) Lookup(name string) (tmpl Template, err error) {
	store := NewCachedStore(s.funcs)
	err = store.Load(s.rootPath)
	if err == nil {
		tmpl, err = store.Lookup(name)
	}
	return tmpl, err
}
