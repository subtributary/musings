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

type ConfigArgs struct {
	BindAddress string
}

type ConfigFile struct {
	// todo: default this to [Und] if not set.
	Locales []localization.Locale `json:"locales"`
}

type Config struct {
	ConfigFile
	ConfigArgs
}

// LoadConfig loads all the configurations needed by the application.
// Configuration sources are "data/config.json", environment variables, and
// command arguments.
//
// If an error occurs, a friendly error message is output to the console.
func LoadConfig() (cfg Config, err error) {
	cfg.ConfigFile, err = loadConfigFile()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
		return
	}

	cfg.ConfigArgs, err = ArgsParser{
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		ProgramName: os.Args[0],
		Getenv:      os.Getenv,
	}.Parse(os.Args[1:])
	return
}

// loadConfigFile loads and parses the project config in "data/config.json".
func loadConfigFile() (ConfigFile, error) {
	path := filepath.Join(DataPath, "config.json")
	contents, err := os.ReadFile(path)
	if err != nil {
		return ConfigFile{}, fmt.Errorf("read config file: %w", err)
	}

	var cfg ConfigFile
	if err = json.Unmarshal(contents, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config file: %w", err)
	}

	return cfg, nil
}

// ArgsParser is for parsing program arguments and environment variables into a
// ConfigArgs. The program arguments are per-run configurations. They do not
// include per-project configurations—those are stored in the config file.
type ArgsParser struct {
	Stdout      io.Writer
	Stderr      io.Writer
	ProgramName string
	Getenv      func(string) string
}

// Parse parses arguments into a ConfigArgs.
// The args passed must start after the program name.
//
// If an error occurs, a friendly error message is output.
func (p ArgsParser) Parse(args []string) (ConfigArgs, error) {
	var cfg ConfigArgs

	cfg.BindAddress = DefaultBindAddress
	if v := p.Getenv("MUSINGS_BIND"); v != "" {
		cfg.BindAddress = v
	}

	var fs flag.FlagSet
	fs.StringVar(&cfg.BindAddress, "bind", cfg.BindAddress, "")
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	err := fs.Parse(args)
	switch {
	case errors.Is(err, flag.ErrHelp):
		p.printUsage()
	case err != nil:
		p.printArgsErr(err)
	case fs.NArg() > 0:
		extraArg := fs.Arg(0)
		err = fmt.Errorf("unexpected positional argument: %s", extraArg)
		p.printArgsErr(err)
	}
	return cfg, err
}

func (p ArgsParser) printArgsErr(err error) {
	_, _ = fmt.Fprintf(p.Stderr, "Error: %v\n", err)
	_, _ = fmt.Fprintf(p.Stderr, "For usage, run: %s --help\n", p.ProgramName)
}

func (p ArgsParser) printUsage() {
	_, _ = fmt.Fprintln(p.Stdout, "Musings website server.")
	_, _ = fmt.Fprintln(p.Stdout)
	_, _ = fmt.Fprintf(p.Stdout, "Usage: %s [options]", p.ProgramName)
	_, _ = fmt.Fprintln(p.Stdout)
	_, _ = fmt.Fprintln(p.Stdout, "Options:")
	_, _ = fmt.Fprintf(p.Stdout, "  --bind <address>  Web endpoint to listen at. [default: %s]\n", DefaultBindAddress)
}
