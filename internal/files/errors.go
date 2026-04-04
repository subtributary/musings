package files

import (
	"errors"
	"fmt"
)

type NotFoundError struct {
	path  string
	inner error
}

func NewNotFoundError(path string, inner error) error {
	return &NotFoundError{
		path:  path,
		inner: inner,
	}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("open path %q", e.path)
}

func (e *NotFoundError) Is(err error) bool {
	var cast *NotFoundError
	return errors.As(err, &cast)
}

func (e *NotFoundError) Unwrap() error {
	return e.inner
}
