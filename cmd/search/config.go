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

	Locale language.Tag
	Query  string
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

func (c *Config) loadAppConfig(args []string) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	localePtr := fs.String("locale", "", "")
	fs.Usage = printUsage

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse args: %w", err)
	}
	if fs.NArg() > 1 {
		return errors.New("too many arguments")
	}

	if *localePtr == "" {
		switch len(c.Global.Locales) {
		case 0:
			// This shouldn't happen, but we can handle it if it does.
			c.Locale = language.Und
		case 1:
			c.Locale = c.Global.Locales[0]
		default:
			return errors.New("locale is ambiguous so must be specified")
		}
	} else if tag, err := language.Parse(*localePtr); err != nil {
		return fmt.Errorf("parse locale: %w", err)
	} else if !slices.Contains(c.Global.Locales, tag) {
		return fmt.Errorf("locale is not enabled")
	} else {
		c.Locale = tag
	}

	c.Query = fs.Arg(0)
	if c.Query == "" {
		return errors.New("query is required")
	}

	return nil
}

func printUsage() {
	program := filepath.Base(os.Args[0])
	fmt.Println("Musings post searcher.")
	fmt.Println()
	fmt.Printf("Usage: %s <query> [options]\n", program)
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --locale <tag>  Set the locale for the index. [default: none]")
}
