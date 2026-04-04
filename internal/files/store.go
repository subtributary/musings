package files

import (
	"os"

	"golang.org/x/text/language"
)

type Store struct {
	matcher  language.Matcher
	rootPath string
}

func NewStore(rootPath string, supportedLocales []string) Store {
	return Store{
		rootPath: rootPath,
	}
}

func (s Store) ReadFile(path string, locale string) ([]byte, error) {
	root, err := os.OpenRoot(s.rootPath)
	if err != nil {
		return nil, NewNotFoundError(s.rootPath, err)
	}
	defer func() { _ = root.Close() }()

	return nil, nil // return locPath.Read(root.FS())
}
