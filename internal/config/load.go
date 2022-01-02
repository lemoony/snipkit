package config

import (
	"os"

	"emperror.dev/errors"
	"github.com/phuslu/log"
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/utils/errorutil"
)

var (
	ErrNoConfigFound = errors.New("no config file use")
	ErrInvalidConfig = errors.New("invalid config file")
	invalidConfig    = Config{}
)

func LoadConfig(v *viper.Viper) (Config, error) {
	if !HasConfig(v) {
		return invalidConfig, errorutil.NewError(ErrNoConfigFound, nil)
	}

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		log.Debug().Str("config file", v.ConfigFileUsed())
	} else {
		return invalidConfig, errorutil.NewError(ErrInvalidConfig, err)
	}

	var wrapper VersionWrapper
	if err := v.Unmarshal(&wrapper); err != nil {
		return invalidConfig, err
	}

	return wrapper.Config, nil
}

func HasConfig(v *viper.Viper) bool {
	_, err := os.Stat(v.ConfigFileUsed())
	return !errors.Is(err, os.ErrNotExist)
}
