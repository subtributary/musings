package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/subtributary/musings/internal/localization"
)

// Global is the global Musings configuration for the build.
// It does not contain per-environment or per-app configurations.
type Global struct {
	Locales []localization.Locale `json:"locales"`
}

// Load loads the global config from the default config file.
func Load() (Global, error) {
	data, err := os.ReadFile(filepath.Join(DataPath, "config.json"))
	if err != nil {
		return Global{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Global
	if err = json.Unmarshal(data, &cfg); err != nil {
		return Global{}, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}

func (g *Global) UnmarshalJSON(data []byte) error {
	type Alias Global

	var state Alias
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	g.Locales = state.Locales
	return nil
}
