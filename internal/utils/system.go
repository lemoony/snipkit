package utils

import (
	"net/url"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
)

const (
	envSnipkitHome = "SNIPKIT_HOME"
)

type System struct {
	Fs afero.Fs

	userDataDir    *string
	userConfigHome *string
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

// WithConfigCome sets the primary config home directory of the user.
func WithConfigCome(configHome string) Option {
	return optionFunc(func(p *System) {
		p.userConfigHome = &configHome
	})
}

// WithUserContainersDir sets the data directory of the user.
func WithUserContainersDir(userContainersDir string) Option {
	return optionFunc(func(p *System) {
		p.userContainersDir = &userContainersDir
	})
}

// WithFS sets the file system.
func WithFS(fs afero.Fs) Option {
	return optionFunc(func(p *System) {
		p.Fs = fs
	})
}

func NewSystem(options ...Option) *System {
	result := System{
		Fs: afero.NewOsFs(),
	}
	for _, option := range options {
		option.apply(&result)
	}
	return &result
}

func NewTestSystem(options ...Option) *System {
	base := afero.NewOsFs()
	roBase := afero.NewReadOnlyFs(base)
	ufs := afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())
	options = append(options, WithFS(ufs))
	return NewSystem(options...)
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

func (s *System) ConfigPath() string {
	return path.Join(s.HomeDir(), "config.yaml")
}

func (s *System) ThemesDir() string {
	return path.Join(s.HomeDir(), "themes/")
}

func (s *System) HomeDir() string {
	dir := os.Getenv(envSnipkitHome)
	if dir != "" {
		return dir
	}

	if s.userConfigHome != nil {
		return *s.userConfigHome
	}

	u, err := url.Parse(xdg.ConfigHome)
	if err != nil {
		panic(err)
	}
	return path.Join(u.Path, "snipkit/")
}
