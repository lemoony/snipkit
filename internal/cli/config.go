package cli

import (
	"os"

	"github.com/phuslu/log"
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/config"
	"github.com/lemoony/snippet-kit/internal/ui/uimsg"
	"github.com/lemoony/snippet-kit/internal/ui/uisurvey"
	"github.com/lemoony/snippet-kit/internal/utils"
)

func ConfigInit(v *viper.Viper) error {
	system, err := utils.NewSystem()
	if err != nil {
		return err
	}

	if _, err := config.LoadConfig(v); err == nil {
		if ok, err2 := uisurvey.ConfirmRecreateConfigFile(v.ConfigFileUsed()); err2 != nil {
			return err2
		} else if !ok {
			log.Info().Msg("User declined to recreate config file")
			return nil
		}
	} else if err == config.ErrNoConfigFound {
		if ok, err2 := uisurvey.ConfirmCreateConfigFile(); err2 != nil {
			return err2
		} else if !ok {
			log.Info().Msg("User declined to create config file")
			return nil
		}
	}

	return config.CreateConfigFile(&system, v)
}

func ConfigEdit(v *viper.Viper) error {
	cfg, err := config.LoadConfig(v)
	if err == config.ErrNoConfigFound {
		uimsg.PrintNoConfig()
		return nil
	}

	return uisurvey.Edit(v.ConfigFileUsed(), cfg.Editor)
}

func ConfigClean(v *viper.Viper) error {
	if !config.HasConfig(v) {
		uimsg.PrintNoConfig()
		return nil
	}

	if ok, err := uisurvey.ConfirmDeleteConfigFile(); err != nil {
		return err
	} else if !ok {
		return nil
	}

	if err := os.RemoveAll(v.ConfigFileUsed()); err != nil {
		return err
	}

	uimsg.PrintConfigDeleted(v.ConfigFileUsed())
	return nil
}
