package config

import (
	"log"

	"golang.org/x/text/language"
)

type Config struct {
	AssetsPath  string // AssetsPath is the path to the website assets directory.
	ContentPath string // ContentPath is the path to the website content directory.
	DataPath    string // DataPath is the path to the data directory.
	BindAddress string // BindAddress is the address for the website to listen on.

	locales string         // Locales is a comma-separated list of locales to consider or support.
	Locales []language.Tag // Locales is a processed version of locales.
}

func (c *Config) process() {
	tags, _, err := language.ParseAcceptLanguage(c.locales)
	if err != nil {
		log.Fatalf("could not parse locales: %v", err)
	}
	if len(tags) == 0 {
		tags = []language.Tag{language.Und}
	}
	c.Locales = tags
}
