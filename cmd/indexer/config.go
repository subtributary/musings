package main

import (
	"flag"
	"os"
)

type Parameter struct {
	name     string
	usage    string
	value    string
	required bool
}

type Config struct {
	presets map[string]string
	custom  map[string]string
}

func (c Config) Preset(key string) string {
	return c.arguments[key]
}

type ConfigBuilder struct {
	//
}

type PresetParameter int

const (
	ContentPath PresetParameter = iota
	DataPath
	Locales
)

type customParameter struct {
	name     string
	usage    string
	value    string
	required bool
}

type ConfigLoader struct {
	presets []PresetParameter
	custom  []customParameter
}

func NewConfigLoader() ConfigLoader {
	return ConfigLoader{}
}

func (c ConfigLoader) Presets(params ...PresetParameter) {
	c.presets = append(c.presets, params...)
}

func (c ConfigLoader) Required(name string, usage string) {
	parameter := customParameter{name, usage, "", true}
	c.custom = append(c.custom, parameter)
}

func (c ConfigLoader) Load(args []string) (Config, error) {
	config := Config{}

	c.loadFromEnv(&config)
	if err := c.loadFromArgs(&config, args); err != nil {
		return Config{}, err
	}

	c.verifyPresets(config)
}

func (c ConfigLoader) loadFromArgs(config *Config, args []string) error {
	cmd := flag.NewFlagSet("config", flag.ExitOnError)

	for _, preset := range c.presets {
		switch preset {
		case ContentPath:
			cmd.StringVar(&config.ContentPath, "content-path", config.ContentPath, "path to content directory")
		case DataPath:
			cmd.StringVar(&config.DataPath, "data-path", config.DataPath, "path to data directory")
		case Locales:
			cmd.StringVar(&config.Locales, "locales", config.Locales, "comma-separated list of locales")
		}
	}

	for _, custom := range c.custom {
		cmd.StringVar(&custom.value, custom.name, "", custom.usage)
	}

	return cmd.Parse(args)
}

func (c ConfigLoader) loadFromEnv(config *Config) {
	for _, preset := range c.presets {
		switch preset {
		case ContentPath:
			config.ContentPath = os.Getenv("CONTENT_PATH")
		case DataPath:
			config.DataPath = os.Getenv("DATA_PATH")
		case Locales:
			config.Locales = os.Getenv("LOCALES")
		}
	}
}

func (c ConfigLoader) verifyPresets(config *Config) {
	for _, preset := range c.presets {
		switch preset {
		case ContentPath:
			if config.ContentPath == "" {
				//
			}
		}
	}
}
