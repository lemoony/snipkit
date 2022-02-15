package system

import (
	"os"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
)

const (
	envSnipkitHome    = "SNIPKIT_HOME"
	fileModeDirectory = os.ModeDir | 0o700
	fileModeConfig    = os.FileMode(0o600)
)

type System struct {
	Fs afero.Fs

	userHome       *string
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

// WithUserHome sets the home directory of the user.
func WithUserHome(userHome string) Option {
	return optionFunc(func(p *System) {
		p.userHome = &userHome
	})
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

func (s *System) HomeEnvValue() string {
	return os.Getenv(envSnipkitHome)
}

func (s *System) UserHome() string {
	if s.userHome != nil {
		return *s.userHome
	}
	return xdg.Home
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

	return path.Join(xdg.ConfigHome, "snipkit/")
}

func (s *System) Remove(path string) {
	if err := s.Fs.Remove(path); err != nil {
		panic(NewErrFileSystem(err, path, "failed to remove"))
	}
}

func (s *System) RemoveAll(path string) {
	if err := s.Fs.RemoveAll(path); err != nil {
		panic(NewErrFileSystem(err, path, "failed to remove all"))
	}
}

func (s *System) DirExists(path string) bool {
	exists, err := afero.DirExists(s.Fs, path)
	if err != nil {
		panic(NewErrFileSystem(err, path, "failed to check if exists"))
	}
	return exists
}

func (s *System) IsEmpty(path string) bool {
	exists, err := afero.IsEmpty(s.Fs, path)
	if err != nil {
		panic(NewErrFileSystem(err, path, "failed to check if empty"))
	}
	return exists
}

// CreatePath returns a suitable location relative to which the file pointed by
// `path` can be written.
func (s *System) CreatePath(path string) {
	dir := filepath.Dir(path)

	if s.DirExists(dir) {
		return
	}

	if err := s.Fs.MkdirAll(dir, fileModeDirectory); err != nil {
		panic(ErrFileSystem{path: dir, msg: "failed to create path", cause: err})
	}
}

func (s *System) WriteFile(path string, data []byte) {
	if err := afero.WriteFile(s.Fs, path, data, fileModeConfig); err != nil {
		panic(ErrFileSystem{path: path, msg: "failed to write file", cause: err})
	}
}

func (s *System) ReadFile(path string) []byte {
	bytes, err := afero.ReadFile(s.Fs, path)
	if err != nil {
		panic(ErrFileSystem{path: path, msg: "failed to read file", cause: err})
	}

	return bytes
}
