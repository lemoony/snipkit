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

// Option configures an App.
type Option interface {
	apply(p *App)
}

// terminalOptionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(a *App)

func (f optionFunc) apply(a *App) {
	f(a)
}

// WithTerminal sets the terminal for the App.
func WithTerminal(t ui.Terminal) Option {
	return optionFunc(func(a *App) {
		a.ui = t
	})
}

type App struct {
	Providers []providers.Provider
	viper     *viper.Viper
	system    *utils.System
	config    *config.Config
	ui        ui.Terminal
}

func NewApp(v *viper.Viper, options ...Option) (*App, error) {
	system := utils.NewSystem()

	app := &App{
		viper:  v,
		system: &system,
		ui:     ui.NewTerminal(),
	}

	for _, o := range options {
		o.apply(app)
	}

	cfg, err := config.LoadConfig(v)
	if err != nil {
		if err == config.ErrNoConfigFound {
			app.ui.PrintError(uimsg.NoConfig())
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
