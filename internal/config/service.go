package config

import (
	"emperror.dev/errors"
	"github.com/phuslu/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/ui"
	"github.com/lemoony/snippet-kit/internal/ui/uimsg"
	"github.com/lemoony/snippet-kit/internal/utils"
	"github.com/lemoony/snippet-kit/internal/utils/errorutil"
)

var (
	ErrNoConfigFound = errors.New("no config file use")
	ErrInvalidConfig = errors.New("invalid config file")
	invalidConfig    = Config{}
)

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
	Create() error
	LoadConfig() (Config, error)
	Edit()
	Clean() error
}

type serviceImpl struct {
	v        *viper.Viper
	system   *utils.System
	terminal ui.Terminal
}

func (s serviceImpl) Create() error {
	if _, err := s.LoadConfig(); err == nil {
		if !s.terminal.Confirm(uimsg.ConfirmRecreateConfigFile(s.v.ConfigFileUsed())) {
			log.Info().Msg("User declined to recreate config file")
			return nil
		}
	} else if err == ErrNoConfigFound && !s.terminal.Confirm(uimsg.ConfirmCreateConfigFile()) {
		log.Info().Msg("User declined to create config file")
		return nil
	}
	return createConfigFile(s.system, s.v, s.terminal)
}

func (s serviceImpl) LoadConfig() (Config, error) {
	if !s.hasConfig() {
		return invalidConfig, errorutil.NewError(ErrNoConfigFound, nil)
	}

	// If a config file is found, read it in.
	if err := s.v.ReadInConfig(); err == nil {
		log.Debug().Str("config file", s.v.ConfigFileUsed())
	} else {
		return invalidConfig, errorutil.NewError(ErrInvalidConfig, err)
	}

	var wrapper VersionWrapper
	if err := s.v.Unmarshal(&wrapper); err != nil {
		return invalidConfig, err
	}

	return wrapper.Config, nil
}

func (s serviceImpl) Edit() {
	cfg, err := s.LoadConfig()
	if err != nil {
		panic(err)
	}

	s.terminal.OpenEditor(s.v.ConfigFileUsed(), cfg.Editor)
}

func (s serviceImpl) Clean() error {
	if !s.hasConfig() {
		s.terminal.PrintError(uimsg.NoConfig())
		return nil
	}

	if !s.terminal.Confirm(uimsg.ConfirmDeleteConfigFile()) {
		return nil
	}

	if err := s.system.Fs.Remove(s.v.ConfigFileUsed()); err != nil {
		return err
	}

	s.terminal.PrintMessage(uimsg.ConfigFileDeleted(s.v.ConfigFileUsed()))
	return nil
}

func (s serviceImpl) hasConfig() bool {
	ok, err := afero.Exists(s.system.Fs, s.v.ConfigFileUsed())
	if err != nil {
		panic(err)
	}
	return ok
}
