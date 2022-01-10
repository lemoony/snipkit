package system

import (
	"fmt"

	"emperror.dev/errors"
)

type ErrFileSystem struct {
	path  string
	msg   string
	cause error
}

func NewErrFileSystem(err error, path, msg string) error {
	return errors.WithStack(ErrFileSystem{path: path, msg: msg, cause: err})
}

func (e ErrFileSystem) Error() string {
	return fmt.Sprintf("%s - %s: %s", e.msg, e.path, e.cause)
}

func (e ErrFileSystem) Is(target error) bool {
	_, ok := target.(ErrFileSystem)
	return ok
}
