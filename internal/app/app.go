package app

import (
	"emperror.dev/errors"
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/system"
)

var ErrNoSnippetsAvailable = errors.New("No snippets are available.")

type App interface {
	LookupSnippet() *model.Snippet
	LookupAndCreatePrintableSnippet() (string, bool)
	LookupAndExecuteSnippet()
	Info()
	AddManager()
	SyncManager()
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

// WithTUI sets the terminal for the App.
func WithTUI(t ui.TUI) Option {
	return optionFunc(func(a *appImpl) {
		a.tui = t
	})
}

// WithProvider sets the provider for the list of manager.
func WithProvider(builder managers.Provider) Option {
	return optionFunc(func(a *appImpl) {
		a.provider = builder
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
		system:   system,
		tui:      ui.NewTUI(),
		provider: managers.NewBuilder(cache.New(system)),
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

	app.tui.ApplyConfig(app.config.Style, system)
	app.managers = app.provider.CreateManager(*app.system, app.config.Manager)

	return app
}

type appImpl struct {
	managers []managers.Manager
	system   *system.System
	config   *config.Config
	tui      ui.TUI

	configService config.Service
	provider      managers.Provider
}

func (a *appImpl) getAllSnippets() []model.Snippet {
	var result []model.Snippet
	for _, manager := range a.managers {
		result = append(result, manager.GetSnippets()...)
	}
	log.Trace().Msgf("Number of available snippets: %d", len(result))
	return result
}
