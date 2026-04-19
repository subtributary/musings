package main

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
)

type Config struct {
	BindAddress         string         // Address to listen at
	ContentPath         string         // Path to website content
	WebPath             string         // Path to website assets
	EnableLiveTemplates bool           // Do not cache template files
	Locales             []language.Tag // Supported locales
}

func (config Config) IsSane() (bool, error) {
	if config.BindAddress == "" {
		return false, errors.New("bind address is not specified")
	}
	if config.ContentPath == "" {
		return false, errors.New("content path is not specified")
	}
	if config.WebPath == "" {
		return false, errors.New("web path is not specified")
	}

	// Errors from inaccessible directories may not be immediate, so we must check their readability here.
	_, err := os.ReadDir(config.ContentPath)
	if err != nil {
		return false, errors.New("content path is not accessible")
	}
	_, err = os.ReadDir(config.WebPath)
	if err != nil {
		return false, errors.New("web path is not accessible")
	}
	_, err = os.ReadDir(config.GetStaticPath())
	if err != nil {
		return false, errors.New("static path is not accessible")
	}
	_, err = os.ReadDir(config.GetTemplatesPath())
	if err != nil {
		return false, errors.New("templates path is not accessible")
	}

	// The format needed for `config.BindAddress` is not well-defined, so it cannot be simply validated.
	// Luckily, an invalid value will cause an immediate error.
	// I will rely on that quick failure instead of trying to validate the endpoint here.

	return true, nil
}

func (config Config) GetStaticPath() string {
	return filepath.Join(config.WebPath, "static")
}

func (config Config) GetTemplatesPath() string {
	return filepath.Join(config.WebPath, "templates")
}
