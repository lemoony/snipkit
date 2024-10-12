package openai

import (
	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/utils/system"
)

// Option configures a Manager.
type Option interface {
	apply(client *Client)
}

// optionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(client *Client)

func (f optionFunc) apply(client *Client) {
	f(client)
}

// WithSystem sets the utils.System instance to be used by Manager.
func WithSystem(system *system.System) Option {
	return optionFunc(func(client *Client) {
		client.system = system
	})
}

func WithConfig(config Config) Option {
	return optionFunc(func(client *Client) {
		client.config = config
	})
}

func WithCache(cache cache.Cache) Option {
	return optionFunc(func(client *Client) {
		client.cache = cache
	})
}
