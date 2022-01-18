package pictarinesnip

import (
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/lemoony/snipkit/internal/utils/system"
)

type Config struct {
	Enabled     bool     `yaml:"enabled" head_comment:"Set to true if you want to use Snip."`
	LibraryPath string   `yaml:"libraryPath" head_comment:"Path to the snippets file."`
	IncludeTags []string `yaml:"includeTags" head_comment:"If this list is not empty, only those snippets that match the listed tags will be provided to you."`
}

func AutoDiscoveryConfig(system *system.System) *Config {
	libPath := findDefaultSnippetsLibrary(system)
	found := libPath != ""

	if libPath == "" {
		libPath = "/path/to/snippets-file"
	}

	return &Config{
		Enabled:     found,
		LibraryPath: libPath,
	}
}

func findDefaultSnippetsLibrary(s *system.System) string {
	containersHome := s.UserContainersHome()
	path := filepath.Join(containersHome, defaultPathContainersLibrary)
	if exists, _ := afero.Exists(s.Fs, path); exists {
		return path
	}
	return ""
}
