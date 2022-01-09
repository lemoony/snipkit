package fslibrary

import (
	"github.com/lemoony/snippet-kit/internal/utils"
)

type Config struct {
	Enabled     bool     `yaml:"enabled" head_comment:"If set to false, the files specified via libraryPath will not be provided to you."`
	LibraryPath []string `yaml:"libraryPath" head_comment:"Paths directories that hold snippets files. Each file must hold one snippet only."`
	SuffixRegex []string `yaml:"suffixRegex" head_comment:"Only files with endings which match one of the listed suffixes will be considered."`
	LazyOpen    bool     `yaml:"lazyOpen" head_comment:"If set to true, the files will not be parsed in advance. This means, only the filename can be used as the snippet name."`
}

func AutoDiscoveryConfig(system *utils.System) Config {
	return Config{
		Enabled:     false,
		LibraryPath: []string{"/path/to/file/system/library"},
		SuffixRegex: []string{".sh"},
		LazyOpen:    false,
	}
}
