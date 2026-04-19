package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/text/language"
)

// IndexConfig is the config needed for indexing
type IndexConfig struct {
	DataPath   string         // Path to data directory
	Locales    []language.Tag // Supported locales
	TargetPath string         // Path to index
}

func (c *IndexConfig) LoadFromEnv() error {
	c.DataPath = os.Getenv("MUSINGS_DATA_PATH")
	c.TargetPath = os.Getenv("MUSINGS_CONTENT_PATH")

	locales, err := parseLocales(os.Getenv("MUSINGS_LOCALES"))
	if err != nil {
		return err
	}
	c.Locales = locales

	return nil
}

// LoadFromArgs loads config from command-line args.
// args should be a slice starting after already-parsed args.
func (c *IndexConfig) LoadFromArgs(args []string) {
	if len(args) > 1 {
		// We should have validated this by now.
		log.Fatal("Arg count for index is wrong.")
	}
	if len(args) == 1 {
		c.TargetPath = args[0]
	}
}

// SearchConfig is the config needed for searching
type SearchConfig struct {
	DataPath string         // Path to data directory
	Locales  []language.Tag // Supported locales
	Query    string
}

func (c *SearchConfig) LoadFromEnv() error {
	c.DataPath = os.Getenv("MUSINGS_DATA_PATH")

	locales, err := parseLocales(os.Getenv("MUSINGS_LOCALES"))
	if err != nil {
		return err
	}
	c.Locales = locales

	return nil
}

// LoadFromArgs loads config from command-line args.
// args should be a slice starting after already-parsed args.
func (c *SearchConfig) LoadFromArgs(args []string) {
	if len(args) != 1 {
		// We should have validated this by now.
		log.Fatal("Arg count for search is wrong.")
	}
	c.Query = args[0]
}

func parseLocales(locales string) ([]language.Tag, error) {
	tags, _, err := language.ParseAcceptLanguage(locales)
	if err != nil {
		return nil, fmt.Errorf("could not parse locales: %w", err)
	}
	if len(tags) == 0 {
		tags = []language.Tag{language.Und}
	}
	return tags, nil
}
