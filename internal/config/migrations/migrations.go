package migrations

import (
	"emperror.dev/errors"
	"gopkg.in/yaml.v3"

	configV110 "github.com/lemoony/snipkit/internal/config/migrations/v110"
	configV111 "github.com/lemoony/snipkit/internal/config/migrations/v111"
	configV112 "github.com/lemoony/snipkit/internal/config/migrations/v112"
)

const (
	Latest = configV112.VersionTo
)

func Migrate(config []byte) []byte {
	for {
		var configMap map[string]interface{}
		if err := yaml.Unmarshal(config, &configMap); err != nil {
			panic(err)
		}

		currentVersion := configMap["version"]

		switch currentVersion {
		case configV110.VersionFrom:
			config = configV110.Migrate(config)
		case configV111.VersionFrom:
			config = configV111.Migrate(config)
		case configV112.VersionFrom:
			config = configV112.Migrate(config)
		case Latest:
			return config
		default:
			panic(errors.Errorf("Unsupported config version: %s", currentVersion))
		}
	}
}
