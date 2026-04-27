package main

import (
	"fmt"
	"os"

	"golang.org/x/text/language"
)

// Config is the global configuration loaded from env.
type Config struct {
	ContentPath string         // Path to content directory
	DataPath    string         // Path to data directory
	Locales     []language.Tag // Supported locales
}

func (c *Config) LoadFromEnv() error {
	c.ContentPath = os.Getenv("MUSINGS_CONTENT_PATH")
	c.DataPath = os.Getenv("MUSINGS_DATA_PATH")

	tags, _, err := language.ParseAcceptLanguage(os.Getenv("MUSINGS_LOCALES"))
	if err != nil {
		return fmt.Errorf("parse locales: %v", err)
	}
	c.Locales = tags

	return nil
}
