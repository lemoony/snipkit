package utils

import (
	"net/url"
	"os/user"
)

type System struct {
	userHomeDir *string
}

// Option configures a Provider.
type Option interface {
	apply(p *System)
}

// optionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(provider *System)

func (f optionFunc) apply(provider *System) {
	f(provider)
}

// WithUserHomeDir sets the home directory of the user.
func WithUserHomeDir(userHomeDir string) Option {
	return optionFunc(func(p *System) {
		p.userHomeDir = &userHomeDir
	})
}

func NewSystem(options ...Option) (System, error) {
	result := System{}
	for _, option := range options {
		option.apply(&result)
	}
	return result, nil
}

func (s *System) UserHomeDir() (string, error) {
	if s.userHomeDir != nil {
		return *s.userHomeDir, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	homeDir, err := url.Parse(currentUser.HomeDir)
	if err != nil {
		return "", err
	}

	return homeDir.Path, nil
}
