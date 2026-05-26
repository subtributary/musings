package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/subtributary/musings/internal/config"
	"golang.org/x/text/language"
)

type Config struct {
	config.Global

	Locales []language.Tag
	Target  string
}

func LoadConfig() (Config, error) {
	base, err := config.Load()
	if err != nil {
		return Config{}, fmt.Errorf("load global config: %w", err)
	}

	cfg := Config{Global: base}
	if err := cfg.loadAppConfig(os.Args[1:]); err != nil {
		return Config{}, fmt.Errorf("load app config: %w", err)
	}

	return cfg, err
}

func (c *Config) LocalizedContentPath(locale language.Tag) string {
	if locale == language.Und {
		return config.ContentPath
	}
	return filepath.Join(config.ContentPath, locale.String())
}

func (c *Config) loadAppConfig(args []string) error {
	fs := flag.NewFlagSet("index", flag.ContinueOnError)
	localePtr := fs.String("locale", "", "")
	fs.Usage = printUsage

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse args: %v", err)
	}
	if fs.NArg() > 1 {
		return errors.New("too many arguments")
	}

	// If locale is not set, default to all locales.
	// If locale is set, validate it then use it as the single locale.
	if *localePtr == "" {
		c.Locales = c.Global.Locales
	} else if tag, err := language.Parse(*localePtr); err != nil {
		return fmt.Errorf("parse locale: %v", err)
	} else if !slices.Contains(c.Global.Locales, tag) {
		return errors.New("locale is not enabled")
	} else {
		c.Locales = []language.Tag{tag}
	}

	c.Target = fs.Arg(0)
	if c.Target != "" && filepath.Ext(c.Target) != ".md" {
		return errors.New("target is not markdown")
	}

	return nil
}

func printUsage() {
	program := filepath.Base(os.Args[0])
	fmt.Println("Musings post indexer.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s [options]\n", program)
	fmt.Printf("  %s [options] <file>\n", program)
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --locale <tag>  Set the locale for the index. [default: all]")
}
