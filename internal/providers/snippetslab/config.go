package snippetslab

import (
	"github.com/lemoony/snippet-kit/internal/utils"
)

type Config struct {
	Enabled     bool
	LibraryPath string
	IncludeTags []string
	ExcludeTags []string
}

func AutoDiscoveryConfig(system *utils.System) Config {
	result := Config{
		Enabled:     false,
		LibraryPath: "/path/to/main.snippetslablibrary",
	}

	libraryURL, err := getLibraryURL(system)
	if err != nil {
		return result
	}

	if ok, err := libraryURL.validate(); err != nil || !ok {
		return result
	} else if basePath, err := libraryURL.basePath(); err == nil {
		result.Enabled = true
		result.LibraryPath = basePath
	}

	return result
}
