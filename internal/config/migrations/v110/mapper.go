package config

import (
	"emperror.dev/errors"
	"gopkg.in/yaml.v3"
)

const (
	VersionFrom = "1.0.0"
	VersionTo   = "1.1.0"
)

func Migrate(old []byte) []byte {
	var config versionWrapper

	if err := yaml.Unmarshal(old, &config); err != nil {
		panic(err)
	}

	if config.Version != VersionFrom {
		panic(errors.Errorf("Invalid version for migration to v1.1.0: %s", config.Version))
	}

	config.Version = VersionTo

	styleCfg := config.Config.Style
	styleCfg["hideKeyMap"] = true

	config.Config.Script = map[string]interface{}{}
	config.Config.Script["shell"] = "/bin/zsh"
	config.Config.Script["parameterMode"] = "SET"
	config.Config.Script["removeComments"] = true

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	return configBytes
}

type versionWrapper struct {
	Version string     `yaml:"version"`
	Config  configV110 `yaml:"config"`
}

type configV110 struct {
	Style              map[string]interface{} `yaml:"style"`
	Editor             string                 `yaml:"editor"`
	DefaultRootCommand string                 `yaml:"defaultRootCommand"`
	Script             map[string]interface{} `yaml:"scripts"`
	Manager            map[string]interface{} `yaml:"manager"`
}
