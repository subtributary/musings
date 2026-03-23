package store

import (
	"fmt"
	"os"

	"github.com/subtributary/musings/internal/localization"
)

type StaticStore struct {
	rootPath string
}

func NewStaticStore(rootPath string) *StaticStore {
	return &StaticStore{
		rootPath: rootPath,
	}
}

// Find searches for the localized variant of a file and returns its resolved path.
//
// If not found, `store.NotFoundError` is returned.
func (s *StaticStore) Find(path string, locale string) (string, error) {
	root, err := os.OpenRoot(s.rootPath)
	if err != nil {
		return "", fmt.Errorf("could not open root: %w", err)
	}
	defer func() { _ = root.Close() }()

	file, err := localization.FindFile(root, path)
	if err != nil {
		return "", newFileNotFoundError(path, err)
	}

	resolvedPath, isFound := file.UseLocale(locale)
	if !isFound {
		return "", newLocaleNotFoundError(path, locale)
	}
	return resolvedPath, nil
}
