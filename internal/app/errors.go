package app

import "errors"

type ArgsError struct {
	Err error
}

func NewArgsErr(err error) error {
	return ArgsError{Err: err}
}

func (e ArgsError) Error() string {
	return e.Error()
}

func AsArgsError(err error) (ArgsError, bool) {
	var cmdErr *ArgsError
	if errors.As(err, &cmdErr) {
		return *cmdErr, true
	}
	return ArgsError{}, false
}
