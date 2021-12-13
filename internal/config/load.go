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

func LoadConfig(viper *viper.Viper) (Config, error) {
	if !HasConfig() {
		return invalidConfig, ErrNoConfigFound
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Str("config file", viper.ConfigFileUsed())
	}

	var wrapper versionWrapper
	if err := viper.Unmarshal(&wrapper); err != nil {
		return invalidConfig, err
	}

	return wrapper.Config, nil
}

func HasConfig() bool {
	_, err := os.Stat(viper.ConfigFileUsed())
	return !errors.Is(err, os.ErrNotExist)
}
