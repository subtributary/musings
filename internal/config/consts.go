package config

import (
	"path/filepath"
	"strings"

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
	localePath := strings.ToLower(locale.String())
	return filepath.Join(ContentPath, localePath)
}
