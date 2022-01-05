package config

import (
	"fmt"

	"emperror.dev/errors"
)

var ErrInvalidConfig = errors.New("invalid config file")

type ErrConfigNotFound struct {
	cfgPath string
}

func (e ErrConfigNotFound) Error() string {
	return fmt.Sprintf("config not found: %s", e.cfgPath)
}

func (e ErrConfigNotFound) Is(target error) bool {
	_, ok := target.(ErrConfigNotFound)
	return ok
}
