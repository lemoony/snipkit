package app

import (
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/config"
	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/providers"
	"github.com/lemoony/snippet-kit/internal/ui"
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

// WithProvidersBuilder sets the builder method for the list of providers.
func WithProvidersBuilder(builder providers.Builder) Option {
	return optionFunc(func(a *App) {
		a.providersBuilder = builder
	})
}

type App struct {
	Providers []providers.Provider
	viper     *viper.Viper
	system    *utils.System
	config    *config.Config
	ui        ui.Terminal

	providersBuilder providers.Builder
}

func NewApp(v *viper.Viper, options ...Option) (*App, error) {
	system := utils.NewSystem()

	app := &App{
		viper:            v,
		system:           &system,
		ui:               ui.NewTerminal(),
		providersBuilder: providers.NewBuilder(),
	}

	for _, o := range options {
		o.apply(app)
	}

	if cfg, err := config.LoadConfig(v); err != nil {
		return nil, err
	} else {
		app.config = &cfg
	}

	app.ui.ApplyConfig(app.config.Style)
	if p, err := app.providersBuilder.BuildProvider(*app.system, app.config.Providers); err != nil {
		return nil, err
	} else {
		app.Providers = p
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
