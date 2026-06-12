package config

import (
	"path/filepath"

	"github.com/subtributary/musings/internal/localization"
)

const (
	ContentPath   = "./content/"
	DataPath      = "./data/"
	StaticPath    = "./web/static/"
	TemplatesPath = "./web/templates/"
)

func LocalizedContentPath(locale localization.Locale) string {
	if locale == localization.UndLocale {
		return ContentPath
	}
	return filepath.Join(ContentPath, locale.Tag)
}
