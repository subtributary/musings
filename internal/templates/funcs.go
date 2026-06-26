package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Funcs struct {
	ContentDir fs.FS
	StaticDir  fs.FS
}

func (f Funcs) ApplyTo(t *template.Template) {
	t.Funcs(template.FuncMap{
		"versioned": f.versioned,
	})
}

func (f Funcs) versioned(assetURL string) (template.URL, error) {
	modified, err := f.modified(assetURL)
	if err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(assetURL)
	if err != nil {
		return "", fmt.Errorf("parse URL: %w", err)
	}

	query := parsedURL.Query()
	query.Set("v", strconv.FormatInt(modified.Unix(), 16))
	parsedURL.RawQuery = query.Encode()

	return template.URL(parsedURL.String()), nil
}

func (f Funcs) modified(name string) (time.Time, error) {
	var dir fs.FS
	if strings.HasPrefix(name, "/_static/") {
		dir = f.StaticDir
		name, _ = strings.CutPrefix(name, "/_static/")
	} else {
		dir = f.ContentDir
		name, _ = strings.CutPrefix(name, "/")
	}

	info, err := fs.Stat(dir, name)
	if err != nil {
		return time.Time{}, fmt.Errorf("stat asset: %w", err)
	}

	return info.ModTime(), nil
}
