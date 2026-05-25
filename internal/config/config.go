package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
)

type configFile struct {
	Locales []string `json:"locales"`
}

func (cf *configFile) apply(c *Global) error {
	// We call a function with a more descriptive name for what we are doing.
	// The caller wants to `apply` the config file to a config;
	// but internally we want to `applyLocales` as the (only) step.
	return cf.applyLocales(c)
}

func (cf *configFile) applyLocales(c *Global) error {
	if len(cf.Locales) == 0 {
		c.Locales = []language.Tag{language.Und}
		return nil
	}

	c.Locales = make([]language.Tag, len(cf.Locales))
	for i, locale := range cf.Locales {
		tag, err := language.Parse(locale)
		if err != nil {
			return fmt.Errorf("parse locale %q: %w", locale, err)
		}
		c.Locales[i] = tag
	}
	return nil
}

// Global is the global Musings configuration for the build.
// It does not contain per-environment or per-app configurations.
type Global struct {
	Locales []language.Tag
}

// Load loads the global config from the default config file.
func Load() (cfg Global, _ error) {
	data, err := os.ReadFile(filepath.Join(DataPath, "config.json"))
	if err != nil {
		return Global{}, fmt.Errorf("read config: %w", err)
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		return Global{}, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}

func (c *Global) UnmarshalJSON(data []byte) error {
	var cf configFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return err
	}
	return cf.apply(c)
}
