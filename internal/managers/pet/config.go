package pet

import (
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/utils/system"
)

type Config struct {
	Enabled      bool     `yaml:"enabled" head_comment:"Set to true if you want to use pet."`
	LibraryPaths []string `yaml:"libraryPaths" head_comment:"List of pet snippet files."`
	IncludeTags  []string `yaml:"includeTags" head_comment:"If this list is not empty, only those snippets that match the listed tags will be provided to you."`
}

func AutoDiscoveryConfig(system *system.System) *Config {
	snippetFilePaths, err := parseSnippetFilePaths(system)
	if err != nil {
		log.Info().Err(err).Msg("failed to discover pet snippet file paths")
	}

	found := len(snippetFilePaths) > 0
	if !found {
		snippetFilePaths = []string{"~/.config/pet/snippet.toml"}
	}

	return &Config{
		Enabled:      found,
		LibraryPaths: snippetFilePaths,
	}
}
