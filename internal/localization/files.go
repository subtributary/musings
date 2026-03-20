package localization

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type LocalizedFile struct {
	variants map[string]string
}

// ScanFiles finds files and returns a map of name to LocalizedFile.
// Files must be named like "{name}.{ext}" or "{name}.{locale}.{ext}".
// The locale, if specified, is lowercased.
func ScanFiles(root *os.Root, path string, ext string) (map[string]LocalizedFile, error) {
	dirFile, err := root.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %w", path, err)
	}
	defer func() { _ = dirFile.Close() }()

	files, err := dirFile.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("could not scan files: %w", err)
	}

	result := make(map[string]LocalizedFile)

	for _, file := range files {
		filename := file.Name()
		if file.IsDir() || filepath.Ext(filename) != ext {
			continue
		}

		var name string
		var locale string

		lastDot := len(filename) - len(ext)
		penultimateDot := strings.LastIndex(filename[:lastDot], ".")
		if penultimateDot != -1 {
			name = filename[:penultimateDot]
			locale = strings.ToLower(filename[penultimateDot+1 : lastDot])
		} else {
			name = filename[:lastDot]
			locale = ""
		}

		if _, ok := result[name]; !ok {
			result[name] = LocalizedFile{
				variants: make(map[string]string),
			}
		}
		result[name].variants[locale] = filepath.Join(path, filename)
	}

	return result, nil
}

// UseLocale finds the best variant for the given locale and returns its file path.
func (f *LocalizedFile) UseLocale(locale string) (string, error) {
	if path, found := FindMatch(f.variants, locale); found {
		return path, nil
	}
	return "", fmt.Errorf("could not resolve for locale %q", locale)
}
