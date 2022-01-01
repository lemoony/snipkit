package cli

import (
	"os"

	"github.com/phuslu/log"
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/config"
	"github.com/lemoony/snippet-kit/internal/ui"
	"github.com/lemoony/snippet-kit/internal/ui/uimsg"
	"github.com/lemoony/snippet-kit/internal/utils"
)

func ConfigInit(v *viper.Viper, term ui.Terminal) error {
	system := utils.NewSystem()

	if _, err := config.LoadConfig(v); err == nil {
		if ok, err2 := term.Confirm(uimsg.ConfirmRecreateConfigFile(v.ConfigFileUsed())); err2 != nil {
			return err2
		} else if !ok {
			log.Info().Msg("User declined to recreate config file")
			return nil
		}
	} else if err == config.ErrNoConfigFound {
		if ok, err2 := term.Confirm(uimsg.ConfirmCreateConfigFile()); err2 != nil {
			return err2
		} else if !ok {
			log.Info().Msg("User declined to create config file")
			return nil
		}
	}

	return config.CreateConfigFile(&system, v, term)
}

func ConfigEdit(v *viper.Viper, term ui.Terminal) error {
	cfg, err := config.LoadConfig(v)
	if err == config.ErrNoConfigFound {
		term.PrintError(uimsg.NoConfig())
		return nil
	}

	return term.OpenEditor(v.ConfigFileUsed(), cfg.Editor)
}

func ConfigClean(v *viper.Viper, term ui.Terminal) error {
	if !config.HasConfig(v) {
		term.PrintError(uimsg.NoConfig())
		return nil
	}

	if ok, err := term.Confirm(uimsg.ConfirmDeleteConfigFile()); err != nil {
		return err
	} else if !ok {
		return nil
	}

	if err := os.RemoveAll(v.ConfigFileUsed()); err != nil {
		return err
	}

	term.PrintMessage(uimsg.ConfigFileDeleted(v.ConfigFileUsed()))
	return nil
}
