package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// findBestVariant find the best-matching filename for the given locale.
// It operates off of an entry of the map returned by scan.
func findBestVariant(variants map[string]string, locale string) (string, error) {
	// todo: get the priorities from elsewhere
	priority := []string{locale, ""}

	for _, p := range priority {
		if filename, ok := variants[p]; ok {
			return filename, nil
		}
	}

	return "", fmt.Errorf("could not resolve variant for locale %s", locale)
}

// scan finds files and returns a map of name and locale to file path.
// Files must be named like "{name}.{ext}" or "{name}.{locale}.{ext}".
// The locale, if specified is lowercased.
func scan(root *os.Root, path string, ext string) (map[string]map[string]string, error) {
	dirFile, err := root.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}
	defer func() { _ = dirFile.Close() }()

	files, err := dirFile.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	result := make(map[string]map[string]string)

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

		if result[name] == nil {
			result[name] = make(map[string]string)
		}
		result[name][locale] = filepath.Join(path, filename)
	}

	return result, nil
}
