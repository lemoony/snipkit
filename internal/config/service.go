package config

import (
	"path/filepath"

	"emperror.dev/errors"
	"github.com/phuslu/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

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
	ConfigFilePath() string
}

type serviceImpl struct {
	v        *viper.Viper
	system   *system.System
	terminal ui.Terminal
}

func (s serviceImpl) Create() {
	_, err := s.LoadConfig()

	create := false
	switch {
	case err == nil:
		create = s.terminal.ConfirmWithHelp(
			uimsg.ConfigFileRecreateConfirm(),
			uimsg.ConfigFileRecreateDescription(s.v.ConfigFileUsed()),
		)
		if !create {
			log.Info().Msg("User declined to recreate config file")
		}
	case errors.Is(err, ErrConfigNotFound{}) && !s.terminal.ConfirmWithHelp(
		uimsg.ConfigFileCreateConfirm(s.v.ConfigFileUsed()),
		uimsg.ConfigFileCreateDescription(s.v.ConfigFileUsed(), s.system.HomeEnvValue()),
	):
		log.Info().Msg("User declined to create config file")
	default:
		create = true
	}

	if create {
		createConfigFile(s.system, s.v, s.terminal)
	}
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

	if s.hasConfig() {
		if s.terminal.Confirm(uimsg.ConfigFileDeleteConfirm(configPath)) {
			s.system.Remove(s.v.ConfigFileUsed())
			s.terminal.PrintMessage(uimsg.ConfigFileDeleted(configPath))
		} else {
			s.terminal.PrintMessage(uimsg.ConfigNotDeleted())
		}
	} else {
		s.terminal.PrintMessage(uimsg.ConfigNotFound(configPath))
	}

	if s.hasThemes() {
		if s.terminal.Confirm(uimsg.ThemesDirDeleteConfirm(s.system.ThemesDir())) {
			s.system.RemoveAll(s.system.ThemesDir())
			s.terminal.PrintMessage(uimsg.ThemesDeleted())
		} else {
			s.terminal.PrintMessage(uimsg.ThemesNotDeleted())
		}
	}

	s.deleteDirectoryIfEmpty(s.system.ThemesDir())
	s.deleteDirectoryIfEmpty(filepath.Dir(s.system.ConfigPath()))

	if exists, _ := afero.DirExists(s.system.Fs, s.system.HomeDir()); exists {
		s.terminal.PrintMessage(uimsg.HomeDirectoryStillExists(s.system.HomeDir()))
	}
}

func (s serviceImpl) ConfigFilePath() string {
	return s.v.ConfigFileUsed()
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
