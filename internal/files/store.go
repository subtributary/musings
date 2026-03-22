package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/subtributary/musings/internal/localization"
)

type Store struct {
	rootPath string
}

func NewStore(filePath string) *Store {
	return &Store{rootPath: filePath}
}

// Find searches for the localization of a file at a route.
// It returns (filePath, nil) if found and ("", nil) if not found.
func (s *Store) Find(route string, locale string) (string, error) {
	root, err := os.OpenRoot(s.rootPath)
	if err != nil {
		return "", fmt.Errorf("could not open content root: %w", err)
	}
	defer func() { _ = root.Close() }()

	path := filepath.Dir(route)
	ext := filepath.Ext(route)
	files, err := localization.ScanFiles(root, path, ext)
	if err != nil {
		return "", fmt.Errorf("could not scan files: %w", err)
	}

	filename := strings.TrimRight(filepath.Base(route), ext)
	if variants, ok := files[filename]; ok {
		if filePath, err := variants.UseLocale(locale); err == nil {
			return filepath.Join(s.rootPath, filePath), nil
		}
	}

	return "", nil
}
