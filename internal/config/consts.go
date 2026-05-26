package config

import (
	"path/filepath"

	"golang.org/x/text/language"
)

const (
	ContentPath   = "./content/"
	DataPath      = "./data/"
	StaticPath    = "./web/static/"
	TemplatesPath = "./web/templates/"
)

func LocalizedContentPath(locale language.Tag) string {
	if locale == language.Und {
		return ContentPath
	}
	return filepath.Join(ContentPath, locale.String())
}
