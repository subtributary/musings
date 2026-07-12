package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/subtributary/musings/internal/localization"
)

const (
	ContentPath = "./content/"
	DataPath    = "./data/"
)

type ConfigFile struct {
	Locales []localization.Locale `json:"Locales"`
}

func LoadConfigFile() (ConfigFile, error) {
	cfg := ConfigFile{}

	configPath := filepath.Join(DataPath, "config.json")
	contents, err := os.ReadFile(configPath)
	if err != nil {
		return ConfigFile{}, fmt.Errorf("read config file: %w", err)
	}

	err = json.Unmarshal(contents, &cfg)
	if err != nil {
		return ConfigFile{}, fmt.Errorf("parse config file: %w", err)
	}

	return cfg, nil
}
