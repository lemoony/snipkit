package filesystem

import (
	"github.com/lemoony/snippet-kit/internal/utils"
)

type Config struct {
	Enabled     bool   `yaml:"enabled" head_comment:"If set to false, the files specified via libraryPath will not be provided to you."`
	LibraryPath string `yaml:"libraryPath" head_comment:"Path to a directory that holds snippets files. Each file must hold one snippet only."`
}

func AutoDiscoveryConfig(system *utils.System) Config {
	return Config{
		Enabled:     false,
		LibraryPath: "/path/to/file/system/library",
	}
}
