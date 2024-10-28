package app

import (
	"time"

	"emperror.dev/errors"
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/system"
)

var ErrNoSnippetsAvailable = errors.New("No snippets are available.")

var ErrSnippetIDNotFound = errors.New("Snippet with ID not found.")

type ErrMigrateConfig struct {
	currentVersion string
	latestVersion  string
}

func (e ErrMigrateConfig) Error() string {
	return uimsg.ConfigNeedsMigration(e.currentVersion, e.latestVersion).RenderWith(style.NoopStyle)
}

func (e ErrMigrateConfig) Is(target error) bool {
	_, ok := target.(ErrMigrateConfig)
	return ok
}

type App interface {
	LookupSnippet() (bool, model.Snippet)
	LookupAndCreatePrintableSnippet() (bool, string)
	LookupSnippetArgs() (bool, string, []model.ParameterValue)
	FindSnippetAndPrint(string, []model.ParameterValue) (bool, string)
	LookupAndExecuteSnippet(bool, bool)
	FindScriptAndExecuteWithParameters(string, []model.ParameterValue, bool, bool)
	ExportSnippets([]ExportField, ExportFormat) string
	GenerateSnippetWithAssistant(string, time.Duration)
	EnableAssistant()
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

// WithAssistantProviderFunc sets the assistant provider.
func WithAssistantProviderFunc(providerFunc func(c assistant.Config) assistant.Assistant) Option {
	return optionFunc(func(a *appImpl) {
		a.assistantProviderFunc = providerFunc
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

// WithCheckNeedsConfigMigration sets if the config file is checked if it is up-to-date.
func WithCheckNeedsConfigMigration(checkNeedsConfigMigration bool) Option {
	return optionFunc(func(a *appImpl) {
		a.checkNeedsConfigMigration = checkNeedsConfigMigration
	})
}

func NewApp(options ...Option) App {
	system := system.NewSystem()

	appCache := cache.New(system)

	app := &appImpl{
		system:   system,
		tui:      ui.NewTUI(),
		provider: managers.NewBuilder(appCache),
		assistantProviderFunc: func(config assistant.Config) assistant.Assistant {
			return assistant.NewBuilder(system, config, appCache)
		},
		cache:                     appCache,
		checkNeedsConfigMigration: true,
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

		if needsMigration, fromVersion := app.configService.NeedsMigration(); needsMigration && app.checkNeedsConfigMigration {
			panic(ErrMigrateConfig{fromVersion, config.Version})
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
	cache    cache.Cache

	configService             config.Service
	provider                  managers.Provider
	assistantProviderFunc     func(assistant.Config) assistant.Assistant
	checkNeedsConfigMigration bool
}

func (a *appImpl) getAllSnippets() []model.Snippet {
	var result []model.Snippet
	for _, manager := range a.managers {
		result = append(result, manager.GetSnippets()...)
	}
	log.Trace().Msgf("Number of available snippets: %d", len(result))
	return result
}
