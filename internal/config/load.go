package config

import (
	"errors"
	"os"

	"github.com/phuslu/log"
	"github.com/spf13/viper"
)

var (
	ErrNoConfigFound = errors.New("no config file use")
	invalidConfig    = Config{}
)

func LoadConfig(v *viper.Viper) (Config, error) {
	if !HasConfig(v) {
		return invalidConfig, ErrNoConfigFound
	}

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		log.Debug().Str("config file", v.ConfigFileUsed())
	} else {
		return invalidConfig, err
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
