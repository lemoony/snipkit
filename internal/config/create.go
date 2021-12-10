package config

import (
	"os"

	"github.com/phuslu/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/lemoony/snippet-kit/internal/providers/filesystem"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
	"github.com/lemoony/snippet-kit/internal/ui/uimsg"
	"github.com/lemoony/snippet-kit/internal/utils"
	"github.com/lemoony/snippet-kit/internal/utils/pathutil"
)

const (
	fileModeConfig = os.FileMode(0o600)
)

func CreateConfigFile(system *utils.System, viper *viper.Viper) error {
	config := versionWrapper{
		Version: "1.0.0",
		Config:  Config{},
	}

	config.Config.Providers.SnippetsLab = snippetslab.AutoDiscoveryConfig(system)
	config.Config.Providers.FileSystem = filesystem.AutoDiscoveryConfig(system)

	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	configPath := viper.ConfigFileUsed()

	log.Debug().Msgf("Going to use config path %s", configPath)
	err = pathutil.CreatePath(configPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, fileModeConfig)
	if err != nil {
		return err
	}

	uimsg.PrintConfigFileCreate(configPath)

	return nil
}
