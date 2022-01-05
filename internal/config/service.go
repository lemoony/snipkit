package config

import (
	"emperror.dev/errors"
	"github.com/phuslu/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/ui"
	"github.com/lemoony/snippet-kit/internal/ui/uimsg"
	"github.com/lemoony/snippet-kit/internal/utils"
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
func WithSystem(system *utils.System) Option {
	return optionFunc(func(s *serviceImpl) {
		s.system = system
	})
}

// NewService creates a new Service.
func NewService(options ...Option) Service {
	service := serviceImpl{
		v:      viper.GetViper(),
		system: utils.NewSystem(),
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
	system   *utils.System
	terminal ui.Terminal
}

func (s serviceImpl) Create() {
	_, err := s.LoadConfig()
	switch {
	case err == nil:
		if !s.terminal.Confirm(uimsg.ConfirmRecreateConfigFile(s.v.ConfigFileUsed())) {
			log.Info().Msg("User declined to recreate config file")
		}
	case errors.Is(err, ErrConfigNotFound{}) && !s.terminal.Confirm(uimsg.ConfirmCreateConfigFile()):
		log.Info().Msg("User declined to create config file")
	default:
		createConfigFile(s.system, s.v, s.terminal)
	}
}

func (s serviceImpl) LoadConfig() (Config, error) {
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
	if !s.hasConfig() {
		panic(ErrConfigNotFound{s.v.ConfigFileUsed()})
	}

	if !s.terminal.Confirm(uimsg.ConfirmDeleteConfigFile()) {
		return
	}

	if err := s.system.Fs.Remove(s.v.ConfigFileUsed()); err != nil {
		panic(errors.Wrapf(err, "failed to remove config file: %s", s.v.ConfigFileUsed()))
	}

	s.terminal.PrintMessage(uimsg.ConfigFileDeleted(s.v.ConfigFileUsed()))
}

func (s serviceImpl) ConfigFilePath() string {
	return s.v.ConfigFileUsed()
}

func (s serviceImpl) hasConfig() bool {
	ok, _ := afero.Exists(s.system.Fs, s.v.ConfigFileUsed())
	return ok
}
