package main

import (
	"errors"
	"fmt"
	"os"
)

type Roots struct {
	Content *os.Root
	Static  *os.Root
}

func OpenRoots() (Roots, error) {
	contentRoot, err := os.OpenRoot(ContentPath)
	if err != nil {
		return Roots{}, fmt.Errorf("open content root: %w", err)
	}

	staticRoot, err := os.OpenRoot(StaticPath)
	if err != nil {
		err1 := fmt.Errorf("open static root: %w", err)
		err2 := contentRoot.Close()
		return Roots{}, errors.Join(err1, err2)
	}

	return Roots{
		Content: contentRoot,
		Static:  staticRoot,
	}, nil
}

func (r Roots) Close() error {
	err1 := r.Content.Close()
	err2 := r.Static.Close()
	return errors.Join(err1, err2)
}
