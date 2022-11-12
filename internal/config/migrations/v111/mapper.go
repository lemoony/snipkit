package config

import (
	"emperror.dev/errors"
	"gopkg.in/yaml.v3"
)

const (
	VersionFrom = "1.1.0"
	VersionTo   = "1.1.1"
)

func Migrate(old []byte) []byte {
	var config versionWrapper

	if err := yaml.Unmarshal(old, &config); err != nil {
		panic(err)
	}

	if config.Version != VersionFrom {
		panic(errors.Errorf("Invalid version for migration to v1.1.1: %s", config.Version))
	}

	config.Version = VersionTo
	config.Config.FuzzySearch = true
	config.Config.Script["execConfirm"] = false
	config.Config.Script["execPrint"] = false

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	return configBytes
}

type versionWrapper struct {
	Version string     `yaml:"version"`
	Config  configV111 `yaml:"config"`
}

type configV111 struct {
	Style              map[string]interface{} `yaml:"style"`
	Editor             string                 `yaml:"editor"`
	DefaultRootCommand string                 `yaml:"defaultRootCommand"`
	FuzzySearch        bool                   `yaml:"fuzzySearch"`
	Script             map[string]interface{} `yaml:"scripts"`
	Manager            map[string]interface{} `yaml:"manager"`
}
