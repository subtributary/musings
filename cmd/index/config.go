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

	TargetLocales []language.Tag
	TargetFile    string
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

	return cfg, nil
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
		return fmt.Errorf("parse args: %w", err)
	}
	if fs.NArg() > 1 {
		return errors.New("too many arguments")
	}

	// If locale is not set, default to all locales.
	// If locale is set, validate it then use it as the single locale.
	if *localePtr == "" {
		c.TargetLocales = c.Global.Locales
	} else if tag, err := language.Parse(*localePtr); err != nil {
		return fmt.Errorf("parse locale: %w", err)
	} else if !slices.Contains(c.Global.Locales, tag) {
		return errors.New("locale is not enabled")
	} else {
		c.TargetLocales = []language.Tag{tag}
	}

	c.TargetFile = fs.Arg(0)
	if c.TargetFile != "" && filepath.Ext(c.TargetFile) != ".md" {
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
