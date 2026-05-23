package templates

type Store interface {
	Lookup(name string) (Template, error)
}
