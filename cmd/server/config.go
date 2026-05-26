package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/subtributary/musings/internal/config"
)

const (
	DefaultBindAddress = ":8080"
)

type Config struct {
	config.Global

	BindAddress   string
	LiveTemplates bool
}

// LoadConfig loads all the configurations needed by the application.
// Configuration sources are "data/config.json", environment variables, and
// command arguments.
//
// This extends `LoadAppConfig` which only loads app-specific configuration.
func LoadConfig() (Config, error) {
	base, err := config.Load()
	if err != nil {
		return Config{}, fmt.Errorf("load global config: %w", err)
	}

	cfg := Config{Global: base}
	if err := cfg.loadAppConfig(os.Args[1:], os.Getenv); err != nil {
		return Config{}, fmt.Errorf("load app config: %w", err)
	}

	return cfg, err
}

// loadAppConfig loads the app-specific configuration from environment
// variables and command arguments. The arguments should not include the
// leading executable name.
func (c *Config) loadAppConfig(args []string, getenv func(string) string) error {
	c.BindAddress = DefaultBindAddress

	if v := getenv("MUSINGS_BIND"); v != "" {
		c.BindAddress = v
	}

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&c.BindAddress, "bind", c.BindAddress, "")
	fs.BoolVar(&c.LiveTemplates, "live-templates", false, "")
	fs.Usage = printUsage

	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("parse args: %v", err)
	} else if fs.NArg() > 0 {
		return errors.New("unexpected positional argument")
	}

	return nil
}

func printUsage() {
	program := filepath.Base(os.Args[0])
	fmt.Println("Musings website server.")
	fmt.Println()
	fmt.Printf("Usage: %s [options]", program)
	fmt.Println()
	fmt.Println("Options:")
	fmt.Printf("  --bind <address>  Web endpoint to listen at. [default: %s]\n", DefaultBindAddress)
	fmt.Println("  --live-templates  Reload templates for every request.")
}
