package config

import (
	"flag"
	"log"
	"os"
	"strings"
)

type Preset int

const (
	unset Preset = iota
	AssetsPath
	ContentPath
	DataPath
	BindAddress
	Locales
)

type parameter struct {
	name     string
	usage    string
	value    string // value is the default.
	required bool
	preset   Preset // preset is the Preset this parameter is for, or unset.
}

type Loader struct {
	params []parameter
}

func NewLoader() Loader {
	return Loader{}
}

// Optional adds an optional command-line parameter to parse.
func (loader *Loader) Optional(name string, usage string, value string) {
	loader.params = append(loader.params, parameter{
		name:  name,
		usage: usage,
		value: value,
	})
}

// Presets adds a collection of preconfigured parameters to the config.
// These are loaded from the environment variables then command-line arguments.
func (loader *Loader) Presets(presets ...Preset) {
	for _, preset := range presets {
		var name, usage string
		switch preset {
		case AssetsPath:
			name = "assets-path"
			usage = "path to assets directory (default: MUSINGS_ASSETS_PATH)"
		case ContentPath:
			name = "content-path"
			usage = "path to content directory (default: MUSINGS_CONTENT_PATH)"
		case DataPath:
			name = "data-path"
			usage = "path to data directory (default: MUSINGS_DATA_PATH)"
		case BindAddress:
			name = "bind"
			usage = "address to bind to (default: MUSINGS_BIND_ADDRESS)"
		case Locales:
			name = "locales"
			usage = "comma-separated list of locales to consider (default: MUSINGS_LOCALES)"
		default:
			panic("invalid preset")
		}

		loader.params = append(loader.params, parameter{
			name:     name,
			usage:    usage,
			required: true,
			preset:   preset,
		})
	}
}

// Required adds a required command-line parameter to parse.
func (loader *Loader) Required(name string, usage string) {
	loader.params = append(loader.params, parameter{
		name:     name,
		usage:    usage,
		required: true,
	})
}

func (loader *Loader) Load(args []string) (*Config, map[string]string) {
	loader.load(args)
	loader.verify()
	return loader.parse()
}

// load parses environment variables then cmd-line arguments for each parameter.
func (loader *Loader) load(args []string) {
	cmd := flag.NewFlagSet("config", flag.ExitOnError)

	for _, param := range loader.params {
		switch param.preset {
		case AssetsPath:
			param.value = os.Getenv("MUSINGS_ASSETS_PATH")
		case ContentPath:
			param.value = os.Getenv("MUSINGS_CONTENT_PATH")
		case DataPath:
			param.value = os.Getenv("MUSINGS_DATA_PATH")
		case BindAddress:
			param.value = os.Getenv("MUSINGS_BIND_ADDRESS")
		case Locales:
			param.value = os.Getenv("MUSINGS_LOCALES")
		}

		cmd.StringVar(&param.value, param.name, param.value, param.usage)
	}

	if err := cmd.Parse(args); err != nil {
		log.Fatalf("error parsing arguments: %v", err)
	}
}

// parse processes the parameters into a Config and an arguments array.
func (loader *Loader) parse() (*Config, map[string]string) {
	config := &Config{}
	arguments := make(map[string]string, len(loader.params))

	for _, param := range loader.params {
		switch param.preset {
		case AssetsPath:
			config.AssetsPath = param.value
		case ContentPath:
			config.ContentPath = param.value
		case DataPath:
			config.DataPath = param.value
		case BindAddress:
			config.BindAddress = param.value
		case Locales:
			config.locales = param.value
		default:
			arguments[param.name] = param.value
		}
	}

	config.process()

	return config, arguments
}

// verify ensures all required parameters are set.
func (loader *Loader) verify() {
	missing := make([]string, 0, len(loader.params))
	for _, param := range loader.params {
		if param.value == "" && param.required {
			missing = append(missing, param.name)
		}
	}

	if len(missing) == 1 {
		log.Fatalf("required parameter not set: %s", missing[0])
	}
	if len(missing) > 1 {
		log.Fatalf("required parameters not set: %s", strings.Join(missing, ", "))
	}
}
