package utils

import (
	"path"

	"github.com/adrg/xdg"
)

type System struct {
	userDataDir    *string
	userConfigDirs *[]string
	// userContainersDir is macOS only
	userContainersDir *string
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

// WithUserDataDir sets the data directory of the user.
func WithUserDataDir(userDataDir string) Option {
	return optionFunc(func(p *System) {
		p.userDataDir = &userDataDir
	})
}

// WithUserConfigDirs sets the config directories of the user.
func WithUserConfigDirs(configDirs []string) Option {
	return optionFunc(func(p *System) {
		p.userConfigDirs = &configDirs
	})
}

// WithUserContainersDir sets the data directory of the user.
func WithUserContainersDir(userContainersDir string) Option {
	return optionFunc(func(p *System) {
		p.userContainersDir = &userContainersDir
	})
}

func NewSystem(options ...Option) (System, error) {
	result := System{}
	for _, option := range options {
		option.apply(&result)
	}
	return result, nil
}

func (s *System) UserDataHome() string {
	if s.userDataDir != nil {
		return *s.userDataDir
	}
	return xdg.DataHome
}

func (s *System) UserConfigDirs() []string {
	if s.userConfigDirs != nil {
		return *s.userConfigDirs
	}
	return xdg.ConfigDirs
}

func (s *System) UserContainersHome() string {
	if s.userContainersDir != nil {
		return *s.userContainersDir
	}
	return path.Join(xdg.Home, "Library/Containers/")
}

func (s *System) UserContainerPreferences(appID string) (string, error) {
	containerDir := s.UserContainersHome()
	return path.Join(containerDir, appID, "Data", "Library", "Preferences"), nil
}
