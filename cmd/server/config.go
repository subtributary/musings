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

func LoadConfig() (Config, error) {
	base, err := config.Load()
	if err != nil {
		return Config{}, err
	}

	cfg, err := LoadAppConfig(os.Args, os.Getenv)
	if err != nil {
		return Config{}, err
	}

	cfg.Global = base
	return cfg, nil
}

func LoadAppConfig(args []string, getenv func(string) string) (cfg Config, err error) {
	cfg.BindAddress = DefaultBindAddress

	if v := getenv("MUSINGS_BIND"); v != "" {
		cfg.BindAddress = v
	}

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&cfg.BindAddress, "bind", cfg.BindAddress, "")
	fs.BoolVar(&cfg.LiveTemplates, "live-templates", false, "")
	fs.Usage = printUsage

	err = fs.Parse(args)
	if err == nil && fs.NArg() > 0 {
		err = errors.New("unexpected positional argument")
	}

	return
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
