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
	c.Locales = make([]language.Tag, len(cf.Locales))

	for i, locale := range cf.Locales {
		tag, err := language.Parse(locale)
		if err != nil {
			return fmt.Errorf("parse locale %q", locale)
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
	path := filepath.Join(DataPath, "config.json")
	if data, err := os.ReadFile(path); err != nil {
		return Global{}, fmt.Errorf("read config: %w", err)
	} else if err = json.Unmarshal(data, &cfg); err != nil {
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
