package app

import (
	"errors"
	"os"
)

type Config struct {
	AssetsDir   string // Directory for website static assets
	ContentDir  string // Directory for website content
	WebEndpoint string // Website endpoint to listen at.
}

func NewConfig() *Config {
	return &Config{}
}

func (config *Config) IsSane() (bool, error) {
	if config.AssetsDir == "" {
		return false, errors.New("assets directory is not specified")
	}
	if config.ContentDir == "" {
		return false, errors.New("content directory is not specified")
	}
	if config.WebEndpoint == "" {
		return false, errors.New("website endpoint is not specified")
	}

	// Errors from inaccessible directories may not be immediate, so we must check their readability here.
	_, err := os.ReadDir(config.AssetsDir)
	if err != nil {
		return false, errors.New("assets directory is not accessible")
	}
	_, err = os.ReadDir(config.ContentDir)
	if err != nil {
		return false, errors.New("content directory is not accessible")
	}

	// The format needed for `config.WebEndpoint` is not well-defined, so it cannot be simply validated.
	// Luckily, an invalid value will cause an immediate error.
	// I will rely on that quick failure instead of trying to validate the endpoint here.

	return true, nil
}
