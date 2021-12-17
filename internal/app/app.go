package app

import (
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/config"
	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/providers"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
	"github.com/lemoony/snippet-kit/internal/ui"
	"github.com/lemoony/snippet-kit/internal/ui/uimsg"
	"github.com/lemoony/snippet-kit/internal/utils"
)

type App struct {
	Providers []providers.Provider
	viper     *viper.Viper
	system    *utils.System
	config    *config.Config
}

func NewApp(v *viper.Viper) (*App, error) {
	system, err := utils.NewSystem()
	if err != nil {
		return nil, err
	}

	app := &App{
		viper:  v,
		system: &system,
	}

	cfg, err := config.LoadConfig(v)
	if err != nil {
		if err == config.ErrNoConfigFound {
			uimsg.PrintNoConfig()
			return nil, nil
		}
		return nil, err
	} else {
		app.config = &cfg
	}

	ui.ApplyConfig(cfg.Style)

	snippetsLab, err := snippetslab.NewProvider(
		snippetslab.WithSystem(&system),
		snippetslab.WithConfig(cfg.Providers.SnippetsLab),
	)
	if err != nil {
		return nil, err
	}

	app.Providers = []providers.Provider{
		snippetsLab,
	}

	return app, nil
}

func (a *App) GetAllSnippets() ([]model.Snippet, error) {
	var result []model.Snippet
	for _, provider := range a.Providers {
		if snippets, err := provider.GetSnippets(); err != nil {
			return nil, err
		} else {
			result = append(result, snippets...)
		}
	}
	return result, nil
}
