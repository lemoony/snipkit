package snippetslab

import (
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Config struct {
	Enabled     bool     `yaml:"enabled" head_comment:"Set to true if you want to use SnippetsLab."`
	LibraryPath string   `yaml:"libraryPath" head_comment:"Path to your *.snippetslablibrary file.\nSnipKit will try to detect this file automatically when generating the config."`
	IncludeTags []string `yaml:"includeTags" head_comment:"If this list is not empty, only those snippets that match the listed tags will be provided to you."`
}

func AutoDiscoveryConfig(system *system.System) *Config {
	result := Config{
		Enabled:     false,
		LibraryPath: "/path/to/main.snippetslablibrary",
	}

	var libraryURL snippetsLabLibrary
	preferencesFilePath, _ := findPreferencesPath(system)
	libraryURL = findLibraryURL(system, preferencesFilePath)

	if ok, err := libraryURL.validate(); err != nil || !ok {
		return &result
	} else if basePath, err := libraryURL.basePath(); err == nil {
		result.Enabled = true
		result.LibraryPath = basePath
	}

	return &result
}
