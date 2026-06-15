package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/subtributary/musings/internal/posts"
)

type Dependencies struct {
	ContentRoot *os.Root
	StaticRoot  *os.Root
	Posts       map[string]*posts.Store // store := Posts[locale.Tag]
	Views       *ViewFactory
}

func LoadDependencies(cfg Config) (d Dependencies, err error) {
	d.ContentRoot, err = os.OpenRoot(ContentPath)
	if err != nil {
		err = fmt.Errorf("open content root: %w", err)
		return d, errors.Join(err, d.Close())
	}

	d.StaticRoot, err = os.OpenRoot(StaticPath)
	if err != nil {
		err = fmt.Errorf("open static root: %w", err)
		return d, errors.Join(err, d.Close())
	}

	d.Views, err = LoadViewFactory(cfg, d.StaticRoot)
	if err != nil {
		err = fmt.Errorf("load views: %w", err)
		return d, errors.Join(err, d.Close())
	}

	d.Posts = make(map[string]*posts.Store)

	// Single "und" posts store if localization is disabled.
	if len(cfg.Locales) == 0 {
		store, err := posts.OpenStore(ContentPath)
		if err != nil {
			err = fmt.Errorf("open posts store: %w", err)
			return d, errors.Join(err, d.Close())
		}
		d.Posts["und"] = store
	}

	// Posts store per configured locale
	for _, locale := range cfg.Locales {
		locPath := filepath.Join(ContentPath, locale.Tag)
		store, err := posts.OpenStore(locPath)
		if err != nil {
			err = fmt.Errorf("open posts store for locale %s: %w", locale.Tag, err)
			return d, errors.Join(err, d.Close())
		}
		d.Posts[locale.Tag] = store
	}

	return d, nil
}

func (d *Dependencies) Close() (err error) {
	if d.ContentRoot != nil {
		err = errors.Join(err, d.ContentRoot.Close())
		d.ContentRoot = nil
	}

	if d.StaticRoot != nil {
		err = errors.Join(err, d.StaticRoot.Close())
		d.StaticRoot = nil
	}

	for loc, store := range d.Posts {
		err = errors.Join(err, store.Close())
		delete(d.Posts, loc)
	}

	return err
}
