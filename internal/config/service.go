package config

import (
	"path/filepath"

	"emperror.dev/errors"
	"github.com/phuslu/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/system"
)

var invalidConfig = Config{}

// Option configures an App.
type Option interface {
	apply(s *serviceImpl)
}

// terminalOptionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(s *serviceImpl)

func (f optionFunc) apply(s *serviceImpl) {
	f(s)
}

// WithTerminal sets the terminal for the Service.
func WithTerminal(t ui.Terminal) Option {
	return optionFunc(func(s *serviceImpl) {
		s.terminal = t
	})
}

// WithViper sets the viper instance for the Service.
func WithViper(v *viper.Viper) Option {
	return optionFunc(func(s *serviceImpl) {
		s.v = v
	})
}

// WithSystem sets the system instance for the Service.
func WithSystem(system *system.System) Option {
	return optionFunc(func(s *serviceImpl) {
		s.system = system
	})
}

// NewService creates a new Service.
func NewService(options ...Option) Service {
	service := serviceImpl{
		v:      viper.GetViper(),
		system: system.NewSystem(),
	}
	for _, o := range options {
		o.apply(&service)
	}
	return service
}

type Service interface {
	Create()
	LoadConfig() (Config, error)
	Edit()
	Clean()
	UpdateManagerConfig(config managers.Config)
	ConfigFilePath() string
}

type serviceImpl struct {
	v        *viper.Viper
	system   *system.System
	terminal ui.Terminal
}

func (s serviceImpl) Create() {
	s.applyConfigTheme()

	recreate := s.hasConfig()
	confirmed := s.terminal.Confirmation(
		uimsg.ConfigFileCreateConfirm(s.v.ConfigFileUsed(), s.system.HomeEnvValue(), recreate),
	)

	if confirmed {
		createConfigFile(s.system, s.v)
	}

	s.terminal.Print(uimsg.ConfigFileCreateResult(confirmed, s.v.ConfigFileUsed(), recreate))
}

func (s serviceImpl) LoadConfig() (Config, error) {
	log.Debug().Msgf("SnipKit Home: %s", s.system.HomeDir())

	if !s.hasConfig() {
		return invalidConfig, ErrConfigNotFound{s.v.ConfigFileUsed()}
	}

	// If a config file is found, read it in.
	if err := s.v.ReadInConfig(); err == nil {
		log.Debug().Str("config file", s.v.ConfigFileUsed())
	} else {
		return invalidConfig, errors.Wrap(ErrInvalidConfig, "failed to read config")
	}

	var wrapper VersionWrapper
	if err := s.v.Unmarshal(&wrapper); err != nil {
		return invalidConfig, err
	}

	return wrapper.Config, nil
}

func (s serviceImpl) Edit() {
	cfgEditor := ""
	if cfg, err := s.LoadConfig(); errors.Is(err, ErrConfigNotFound{}) {
		panic(err)
	} else {
		cfgEditor = cfg.Editor
	}

	s.terminal.OpenEditor(s.v.ConfigFileUsed(), cfgEditor)
}

func (s serviceImpl) Clean() {
	configPath := s.v.ConfigFileUsed()
	s.applyConfigTheme()

	if s.hasConfig() {
		confirmed := s.terminal.Confirmation(uimsg.ConfigFileDeleteConfirm(configPath))
		if confirmed {
			s.system.Remove(s.v.ConfigFileUsed())
		}
		s.terminal.Print(uimsg.ConfigFileDeleteResult(confirmed, s.v.ConfigFileUsed()))
	} else {
		s.terminal.Print(uimsg.ConfigNotFound(configPath))
	}

	if s.hasThemes() {
		confirmed := s.terminal.Confirmation(uimsg.ThemesDeleteConfirm(s.system.ThemesDir()))
		if confirmed {
			s.system.RemoveAll(s.system.ThemesDir())
		}
		s.terminal.Print(uimsg.ThemesDeleteResult(confirmed, s.system.ThemesDir()))
	}

	s.deleteDirectoryIfEmpty(s.system.ThemesDir())
	s.deleteDirectoryIfEmpty(filepath.Dir(s.system.ConfigPath()))

	if exists, _ := afero.DirExists(s.system.Fs, s.system.HomeDir()); exists {
		s.terminal.Print(uimsg.HomeDirectoryStillExists(s.system.HomeDir()))
	}
}

func (s serviceImpl) ConfigFilePath() string {
	return s.v.ConfigFileUsed()
}

func (s serviceImpl) UpdateManagerConfig(managerConfig managers.Config) {
	config, err := s.LoadConfig()
	if err != nil {
		panic(errors.Wrapf(ErrInvalidConfig, "failed to load config: %s", err.Error()))
	}

	if cfg := managerConfig.FsLibrary; cfg != nil {
		config.Manager.FsLibrary = cfg
	}
	if cfg := managerConfig.SnippetsLab; cfg != nil {
		config.Manager.SnippetsLab = cfg
	}
	if cfg := managerConfig.PictarineSnip; cfg != nil {
		config.Manager.PictarineSnip = cfg
	}

	bytes := SerializeToYamlWithComment(wrap(config))
	s.system.WriteFile(s.ConfigFilePath(), bytes)
}

func (s serviceImpl) hasConfig() bool {
	ok, _ := afero.Exists(s.system.Fs, s.v.ConfigFileUsed())
	return ok
}

func (s serviceImpl) hasThemes() bool {
	themesDir := s.system.ThemesDir()
	if exists, _ := afero.DirExists(s.system.Fs, themesDir); !exists {
		return false
	}
	return !s.system.IsEmpty(themesDir)
}

func (s serviceImpl) deleteDirectoryIfEmpty(path string) {
	if s.system.DirExists(path) && s.system.IsEmpty(path) {
		s.system.Remove(path)
	}
}

func (s serviceImpl) applyConfigTheme() {
	cfg, err := s.LoadConfig()
	if err == nil {
		ui.ApplyConfig(cfg.Style, s.system)
	} else {
		ui.ApplyConfig(ui.DefaultConfig(), s.system)
	}
}
