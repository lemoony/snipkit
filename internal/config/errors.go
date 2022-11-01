package config

import (
	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
)

var ErrInvalidConfig = errors.New("invalid config file")

type ErrConfigNotFound struct {
	cfgPath string
}

func (e ErrConfigNotFound) Error() string {
	return uimsg.ConfigNotFound(e.cfgPath).RenderWith(style.NoopStyle)
}

func (e ErrConfigNotFound) Is(target error) bool {
	_, ok := target.(ErrConfigNotFound)
	return ok
}

type ErrMigrateConfig struct {
	currentVersion string
	latestVersion  string
}

func (e ErrMigrateConfig) Error() string {
	return uimsg.ConfigNeedsMigration(e.currentVersion, e.latestVersion).RenderWith(style.NoopStyle)
}

func (e ErrMigrateConfig) Is(target error) bool {
	_, ok := target.(ErrMigrateConfig)
	return ok
}
