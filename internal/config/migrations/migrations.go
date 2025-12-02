package migrations

import (
	"emperror.dev/errors"
	"gopkg.in/yaml.v3"

	configV110 "github.com/lemoony/snipkit/internal/config/migrations/v110"
	configV111 "github.com/lemoony/snipkit/internal/config/migrations/v111"
	configV120 "github.com/lemoony/snipkit/internal/config/migrations/v120"
	configV130 "github.com/lemoony/snipkit/internal/config/migrations/v130"
)

const (
	Latest = configV130.VersionTo
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
		case configV120.VersionFrom:
			config = configV120.Migrate(config)
		case configV130.VersionFrom:
			config = configV130.Migrate(config)
		case Latest:
			return config
		default:
			panic(errors.Errorf("Unsupported config version: %s", currentVersion))
		}
	}
}
