package templates

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/url"
	"strings"
)

type Funcs struct {
	StaticDir fs.FS
}

func (f Funcs) ApplyTo(t *template.Template) {
	t.Funcs(template.FuncMap{
		"versioned": f.versioned,
	})
}

func (f Funcs) versioned(assetURL string) (template.URL, error) {
	const timestamp = "20060102150405"

	if f.StaticDir == nil {
		return "", errors.New("static dir is not defined")
	}

	// Get the modified time from the file.
	assetPath, ok := strings.CutPrefix(assetURL, "/_static/")
	if !ok {
		return "", errors.New("URL is not to a static asset")
	}
	info, err := fs.Stat(f.StaticDir, assetPath)
	if err != nil {
		return "", fmt.Errorf("stat asset: %w", err)
	}
	modified := info.ModTime()

	// Append the modified time to the URL.
	parsedURL, err := url.Parse(assetURL)
	if err != nil {
		return "", fmt.Errorf("parse URL: %w", err)
	}

	query := parsedURL.Query()
	query.Set("v", modified.UTC().Format(timestamp))
	parsedURL.RawQuery = query.Encode()

	return template.URL(parsedURL.String()), nil
}
