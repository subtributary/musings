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

// FindFile finds a specific file and returns its localization options.
// Files must be named like "{name}.{ext}" or "{name}.{locale}.{ext}".
// The locale, if specified, is lowercased.
func FindFile(root *os.Root, path string) (result LocalizedFile, err error) {
	filename := filepath.Base(path)

	// We can't ask the system to filter by file prefix,
	// so we have to manually enumerate the files in the directory.
	// We could optimize this by doing the enumeration and filtering here,
	// but I feel like the benefit is not worth the uglier code.
	files, err := ScanFiles(root, filepath.Dir(path))
	if err != nil {
		return
	} else if found, ok := files[filename]; ok {
		result = found
	} else {
		err = fmt.Errorf("file %q not found", filename)
	}

	return
}

// ScanFiles finds files and returns a map of filename to LocalizedFile.
// Files must be named like "{name}.{ext}" or "{name}.{locale}.{ext}".
// The locale, if specified, is lowercased.
func ScanFiles(root *os.Root, path string) (map[string]LocalizedFile, error) {
	dirFile, err := root.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %w", path, err)
	}
	defer func() { _ = dirFile.Close() }()

	files, err := dirFile.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("could not read dir %q: %w", path, err)
	}

	result := make(map[string]LocalizedFile)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		name, locale := parseFilename(filename)

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
func (f *LocalizedFile) UseLocale(locale string) (string, bool) {
	return FindMatch(f.variants, locale)
}

// parseFilename parses a filename into name and locale per our conventions:
//   - Filename "{name}.{ext}" becomes (filename, "")
//   - Filename "{name}.{locale}.{ext}" becomes ("{name}.{ext}", locale)
//
// The locale, is specified, is lowercased.
func parseFilename(filename string) (name string, locale string) {
	ext := filepath.Ext(filename)

	lastDot := len(filename) - len(ext)
	penultimateDot := strings.LastIndex(filename[:lastDot], ".")
	if penultimateDot != -1 {
		name = filename[:penultimateDot] + ext
		locale = strings.ToLower(filename[penultimateDot+1 : lastDot])
	} else {
		name = filename
	}

	return
}
