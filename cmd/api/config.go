package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/subtributary/musings/internal/app"
)

const DefaultBindAddress = ":8081"

type Config struct {
	app.ConfigFile
	BindAddress string
}

// LoadConfig loads all the configurations needed by the application.
// Configuration sources are "data/config.json" and command arguments.
func LoadConfig() (Config, error) {
	configFile, err := app.LoadConfigFile()
	if err != nil {
		return Config{}, fmt.Errorf("load config file: %w", err)
	}

	cfg := Config{
		ConfigFile:  configFile,
		BindAddress: DefaultBindAddress,
	}

	var fs flag.FlagSet
	fs.StringVar(&cfg.BindAddress, "bind", cfg.BindAddress, "")
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	err = fs.Parse(os.Args[1:])
	if err != nil {
		return Config{}, app.NewArgsErr(err)
	}
	if fs.NArg() > 0 {
		err := fmt.Errorf("unexpected positional argument: %s", fs.Arg(0))
		return Config{}, app.NewArgsErr(err)
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
	_, _ = fmt.Fprintln(os.Stdout, "Musings website CMS.")
	_, _ = fmt.Fprintln(os.Stdout)
	_, _ = fmt.Fprintf(os.Stdout, "Usage: %s [options]", programName)
	_, _ = fmt.Fprintln(os.Stdout)
	_, _ = fmt.Fprintln(os.Stdout, "Options:")
	_, _ = fmt.Fprintf(os.Stdout, "  --bind <address>  Web endpoint to listen at. [default: %s]\n", DefaultBindAddress)
}
