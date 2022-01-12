package app

import (
	"emperror.dev/errors"
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/providers"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/system"
)

var ErrNoSnippetsAvailable = errors.New("No snippets are available.")

type App interface {
	LookupSnippet() *model.Snippet
	LookupAndCreatePrintableSnippet() (string, bool)
	LookupAndExecuteSnippet()
	Info()
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

func NewApp(options ...Option) App {
	system := system.NewSystem()

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
			panic(errors.WithStack(err))
		} else {
			app.config = &cfg
		}
	}

	if app.config == nil {
		panic("no config provided")
	}

	app.ui.ApplyConfig(app.config.Style, system)
	if p, err := app.providersBuilder.BuildProvider(*app.system, app.config.Providers); err != nil {
		panic(err)
	} else {
		app.Providers = p
	}

	return app
}

type appImpl struct {
	Providers []providers.Provider
	system    *system.System
	config    *config.Config
	ui        ui.Terminal

	configService    config.Service
	providersBuilder providers.Builder
}

func (a *appImpl) getAllSnippets() []model.Snippet {
	var result []model.Snippet
	for _, provider := range a.Providers {
		result = append(result, provider.GetSnippets()...)
	}
	log.Trace().Msgf("Number of available snippets: %d", len(result))
	return result
}
