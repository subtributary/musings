package config

import (
	"path/filepath"
	"strings"

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
	localePath := strings.ToLower(locale.Tag)
	return filepath.Join(ContentPath, localePath)
}
