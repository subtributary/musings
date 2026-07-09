package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/subtributary/musings/internal/localization"
)

const (
	ContentPath        = "./content/"
	DataPath           = "./data/"
	DefaultBindAddress = ":8080"
)

type ArgsError struct {
	Err error
}

func (e ArgsError) Error() string {
	return e.Error()
}

func ToArgsError(err error) (ArgsError, bool) {
	var cmdErr *ArgsError
	if errors.As(err, &cmdErr) {
		return *cmdErr, true
	}
	return ArgsError{}, false
}

type Config struct {
	BindAddress string
	Locales     []localization.Locale `json:"locales"`
}

// LoadConfig loads all the configurations needed by the application.
// Configuration sources are "data/config.json", environment variables, and
// command arguments.
//
// If an error occurs, a friendly error message is output to the console.
func LoadConfig() (cfg Config, err error) {
	cfg.BindAddress = DefaultBindAddress

	configPath := filepath.Join(DataPath, "config.json")
	if contents, err := os.ReadFile(configPath); err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	} else if err := json.Unmarshal(contents, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config file: %w", err)
	}

	if value := os.Getenv("MUSINGS_BIND"); value != "" {
		cfg.BindAddress = value
	}

	var fs flag.FlagSet
	fs.StringVar(&cfg.BindAddress, "bind", cfg.BindAddress, "")
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	if err := fs.Parse(os.Args[1:]); err != nil {
		return Config{}, ArgsError{Err: err}
	} else if fs.NArg() > 0 {
		err := fmt.Errorf("unexpected positional argument: %s", fs.Arg(0))
		return Config{}, ArgsError{Err: err}
	}

	return cfg, nil
}

func PrintArgsErr(err error) {
	programName := os.Args[0]
	_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	_, _ = fmt.Fprintf(os.Stderr, "For usage, run: %s --help\n", programName)
}

func PrintUsage() {
	programName := os.Args[0]
	_, _ = fmt.Fprintln(os.Stdout, "Musings website server.")
	_, _ = fmt.Fprintln(os.Stdout)
	_, _ = fmt.Fprintf(os.Stdout, "Usage: %s [options]", programName)
	_, _ = fmt.Fprintln(os.Stdout)
	_, _ = fmt.Fprintln(os.Stdout, "Options:")
	_, _ = fmt.Fprintf(os.Stdout, "  --bind <address>  Web endpoint to listen at. [default: %s]\n", DefaultBindAddress)
}
