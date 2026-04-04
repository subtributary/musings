package files

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
)

type DirEntry struct {
	name    string
	wrapped fs.DirEntry
}

// newDirEntry creates a new DirEntry with a name override.
func newDirEntry(wrapped fs.DirEntry, name string) DirEntry {
	return DirEntry{
		name:    name,
		wrapped: wrapped,
	}
}

func (d DirEntry) Name() string {
	return d.name
}

func (d DirEntry) IsDir() bool {
	return d.wrapped.IsDir()
}

func (d DirEntry) Type() fs.FileMode {
	return d.wrapped.Type()
}

func (d DirEntry) Info() (fs.FileInfo, error) {
	return d.wrapped.Info()
}

// LocalizedFS is a file system that supports localized versions of files.
//
// Localized versions of files should be named with the lowercased canonical
// locale string as a secondary extension in front of the primary extension;
// for example, "hello.en.md" is the English version of "hello.md".
type LocalizedFS struct {
	tag     language.Tag
	wrapped fs.ReadDirFS
}

func NewLocalizedFS(wrapped fs.ReadDirFS, tag language.Tag) LocalizedFS {
	return LocalizedFS{tag, wrapped}
}

func (f LocalizedFS) Open(name string) (fs.File, error) {
	if original, _ := parseFileName(name); original != name {
		// To avoid ambiguity, ignore requests that include a locale.
		return nil, fs.ErrNotExist
	}

	// Attempt to read localized variants
	ext := filepath.Ext(name)
	prefix := strings.TrimSuffix(name, ext)
	for tag := f.tag; !tag.IsRoot(); tag = tag.Parent() {
		localeExt := "." + strings.ToLower(tag.String())
		path := prefix + localeExt + ext
		file, err := f.wrapped.Open(path)
		if err == nil {
			return file, nil
		}
	}

	// Fallback to original
	return f.wrapped.Open(name)
}

// ReadDir lists all files supported by the configured locale.
//
// If files for more than one locale are desired, then Scan is more performant than
// calling ReadDir for multiple localized file systems.
func (f LocalizedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	files, err := f.wrapped.ReadDir(name)
	if err != nil {
		return nil, err
	}

	groupedFiles := scan(files, []language.Tag{f.tag})
	return groupedFiles[f.tag], nil
}

// Scan returns the files in a directory with the locale stripped from the filename
// and used to group files together by locale. Files with an unsupported locale are
// not listed.
func Scan(dir fs.FS, tags []language.Tag) (map[language.Tag][]fs.DirEntry, error) {
	files, err := fs.ReadDir(dir, ".")
	if err != nil {
		return nil, err
	}
	return scan(files, tags), nil
}

func scan(files []fs.DirEntry, tags []language.Tag) map[language.Tag][]fs.DirEntry {
	result := make(map[language.Tag][]fs.DirEntry)
	visited := make(map[language.Tag]map[string]struct{})
	for _, tag := range tags {
		result[tag] = make([]fs.DirEntry, 0)
		visited[tag] = make(map[string]struct{})
	}

	for _, file := range files {
		// Directories are not localized.
		if file.IsDir() {
			for _, tag := range tags {
				result[tag] = append(result[tag], file)
			}
			continue
		}

		fileName, fileTag := parseFileName(file.Name())
		dirEntry := newDirEntry(file, fileName)

		// Unlocalized filenames are added to all sets.
		if fileTag == language.Und {
			for _, tag := range tags {
				if _, ok := visited[tag][fileName]; ok {
					continue
				}
				result[tag] = append(result[tag], dirEntry)
				visited[tag][fileName] = struct{}{}
			}
			continue
		}

		// Add to all sets that understand the file tag.
		found := false
		for _, tag := range tags {
			if _, ok := visited[tag][fileName]; ok {
				continue
			}
			if language.Comprehends(tag, fileTag) == language.No {
				continue
			}
			result[tag] = append(result[tag], dirEntry)
			visited[tag][fileName] = struct{}{}
			found = true
		}

		// Helpful log message
		if !found {
			unlocalizedName := file.Name()
			log.Printf("file skipped because of disabled locale: %s", unlocalizedName)
		}
	}

	return result
}

// parseFileName extracts the locale from the filename and returns it as a
// language.Tag alongside the unlocalized filename.
func parseFileName(name string) (string, language.Tag) {
	ext := filepath.Ext(name)
	prefix := name[:len(name)-len(ext)]
	localeExt := filepath.Ext(prefix)

	if localeExt == "" {
		return name, language.Und
	}

	tag, err := language.Parse(localeExt[1:])
	if err != nil {
		// Treat as no filename localization but with a warning.
		log.Printf("invalid locale in filename: %s", name)
		return name, language.Und
	}

	name = prefix[:len(prefix)-len(localeExt)] + ext
	return name, tag
}
