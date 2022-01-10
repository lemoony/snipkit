package config

import (
	"emperror.dev/errors"

	"github.com/lemoony/snippet-kit/internal/ui/uimsg"
)

var ErrInvalidConfig = errors.New("invalid config file")

type ErrConfigNotFound struct {
	cfgPath string
}

func (e ErrConfigNotFound) Error() string {
	return uimsg.ConfigNotFound(e.cfgPath)
}

func (e ErrConfigNotFound) Is(target error) bool {
	_, ok := target.(ErrConfigNotFound)
	return ok
}
