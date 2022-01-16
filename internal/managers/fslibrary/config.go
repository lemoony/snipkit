package fslibrary

import (
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Config struct {
	Enabled            bool     `yaml:"enabled" head_comment:"If set to false, the files specified via libraryPath will not be provided to you."`
	LibraryPath        []string `yaml:"libraryPath" head_comment:"Paths directories that hold snippets files. Each file must hold one snippet only."`
	SuffixRegex        []string `yaml:"suffixRegex" head_comment:"Only files with endings which match one of the listed suffixes will be considered."`
	LazyOpen           bool     `yaml:"lazyOpen" head_comment:"If set to true, the files will not be parsed in advance. This means, only the filename can be used as the snippet name."`
	HideTitleInPreview bool     `yaml:"hideTitleInPreview" head_comment:"If set to true, the title comment will not be shown in the preview window."`
}

func AutoDiscoveryConfig(system *system.System) Config {
	return Config{
		Enabled:            false,
		LibraryPath:        []string{"/path/to/file/system/library"},
		SuffixRegex:        []string{".sh"},
		LazyOpen:           false,
		HideTitleInPreview: false,
	}
}
