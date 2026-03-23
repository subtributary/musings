package store

import (
	"errors"
	"fmt"
)

type NotFoundError struct {
	inner   error
	message string
}

func (e *NotFoundError) Error() string {
	return e.message
}

func (e *NotFoundError) Unwrap() error {
	return e.inner
}

func (e *NotFoundError) Is(err error) bool {
	var cast *NotFoundError
	return errors.As(err, &cast)
}

func newFileNotFoundError(path string, inner error) error {
	return &NotFoundError{
		inner:   inner,
		message: fmt.Sprintf("file not found: %s", path),
	}
}

func newLocaleNotFoundError(path string, locale string) error {
	return &NotFoundError{
		message: fmt.Sprintf("locale %q not found for file: %s", locale, path),
	}
}
