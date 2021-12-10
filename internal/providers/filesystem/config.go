package filesystem

import (
	"github.com/lemoony/snippet-kit/internal/utils"
)

type Config struct {
	Enabled     bool
	LibraryPath string
}

func AutoDiscoveryConfig(system *utils.System) Config {
	return Config{
		Enabled:     false,
		LibraryPath: "/path/to/file/system/library",
	}
}
