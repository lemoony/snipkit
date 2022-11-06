package migrations

import (
	"emperror.dev/errors"
	"gopkg.in/yaml.v3"

	configV111 "github.com/lemoony/snipkit/internal/config/migrations/v111"
)

const (
	Latest = configV111.VersionTo
)

func Migrate(config []byte) []byte {
	for {
		var configMap map[string]interface{}
		if err := yaml.Unmarshal(config, &configMap); err != nil {
			panic(err)
		}

		currentVersion := configMap["version"]

		switch currentVersion {
		case configV111.VersionFrom:
			config = configV111.Migrate(config)
		case Latest:
			return config
		default:
			panic(errors.Errorf("Unsupported config version: %s", currentVersion))
		}
	}
}
