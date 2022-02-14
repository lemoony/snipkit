package githubgist

import (
	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/utils/system"
)

// Option configures a Manager.
type Option interface {
	apply(p *Manager)
}

// optionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(manager *Manager)

func (f optionFunc) apply(manager *Manager) {
	f(manager)
}

// WithSystem sets the utils.System instance to be used by Manager.
func WithSystem(system *system.System) Option {
	return optionFunc(func(p *Manager) {
		p.system = system
	})
}

func WithConfig(config Config) Option {
	return optionFunc(func(p *Manager) {
		p.config = config
	})
}

func WithCache(cache cache.Cache) Option {
	return optionFunc(func(p *Manager) {
		p.cache = cache
	})
}
