package app

import (
	"emperror.dev/errors"

	"github.com/lemoony/snippet-kit/internal/config"
	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/providers"
	"github.com/lemoony/snippet-kit/internal/ui"
	"github.com/lemoony/snippet-kit/internal/utils"
)

type App interface {
	LookupSnippet() (*model.Snippet, error)
	LookupAndCreatePrintableSnippet() (string, error)
	LookupAndExecuteSnippet() error
	Info() error
}

// Option configures an App.
type Option interface {
	apply(p *appImpl)
}

// terminalOptionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(a *appImpl)

func (f optionFunc) apply(a *appImpl) {
	f(a)
}

// WithTerminal sets the terminal for the App.
func WithTerminal(t ui.Terminal) Option {
	return optionFunc(func(a *appImpl) {
		a.ui = t
	})
}

// WithProvidersBuilder sets the builder method for the list of providers.
func WithProvidersBuilder(builder providers.Builder) Option {
	return optionFunc(func(a *appImpl) {
		a.providersBuilder = builder
	})
}

// WithConfig sets the config for the App.
func WithConfig(config config.Config) Option {
	return optionFunc(func(a *appImpl) {
		a.config = &config
	})
}

// WithConfigService sets the config service for the App.
func WithConfigService(service config.Service) Option {
	return optionFunc(func(a *appImpl) {
		a.configService = service
	})
}

func NewApp(options ...Option) (App, error) {
	system := utils.NewSystem()

	app := &appImpl{
		system:           system,
		ui:               ui.NewTerminal(),
		providersBuilder: providers.NewBuilder(),
	}

	for _, o := range options {
		o.apply(app)
	}

	if app.configService != nil {
		if cfg, err := app.configService.LoadConfig(); err != nil {
			return nil, err
		} else {
			app.config = &cfg
		}
	}

	if app.config == nil {
		return nil, errors.New("no config provided")
	}

	app.ui.ApplyConfig(app.config.Style)
	if p, err := app.providersBuilder.BuildProvider(*app.system, app.config.Providers); err != nil {
		return nil, err
	} else {
		app.Providers = p
	}

	return app, nil
}

type appImpl struct {
	Providers []providers.Provider
	system    *utils.System
	config    *config.Config
	ui        ui.Terminal

	configService    config.Service
	providersBuilder providers.Builder
}

func (a *appImpl) getAllSnippets() ([]model.Snippet, error) {
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
