package config

import (
	"github.com/spf13/viper"

	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/system"
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

// WithTerminal sets the tui for the Service.
func WithTerminal(t ui.TUI) Option {
	return optionFunc(func(s *serviceImpl) {
		s.tui = t
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
