package templates

type LiveStore struct {
	rootPath string
}

func NewLiveStore(rootPath string) LiveStore {
	return LiveStore{rootPath: rootPath}
}

func (s LiveStore) Lookup(name string) (tmpl Template, err error) {
	store := NewCachedStore()
	err = store.Load(s.rootPath)
	if err == nil {
		tmpl, err = store.Lookup(name)
	}
	return tmpl, err
}
