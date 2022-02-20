package masscode

import (
	"path/filepath"

	"github.com/lemoony/snipkit/internal/utils/system"
)

type Version string

const (
	version1 = Version("v1")
	version2 = Version("v2")
)

type Config struct {
	Enabled      bool     `yaml:"enabled" head_comment:"Set to true if you want to use pet."`
	MassCodeHome string   `yaml:"massCodeHome" head_comment:"Path to the massCode directory containing the db files."`
	Version      Version  `version:"Version of massCode. Allowed values: v1, v2."`
	IncludeTags  []string `yaml:"includeTags" head_comment:"If this list is not empty, only those Snippets that match the listed Tags will be provided to you."`
}

type autoDetectResult struct {
	found   bool
	version Version
	path    string
}

func AutoDiscoveryConfig(system *system.System) *Config {
	autoDetect := findMassCodeHome(system)
	return &Config{
		Enabled:      autoDetect.found,
		Version:      autoDetect.version,
		MassCodeHome: autoDetect.path,
		IncludeTags:  []string{},
	}
}

func findMassCodeHome(sys *system.System) autoDetectResult {
	defaultHome := filepath.Join(sys.UserHome(), defaultMassCodeHomePath)

	result := autoDetectResult{
		found:   false,
		version: version1,
		path:    defaultHome,
	}

	if v2DBFile := filepath.Join(defaultHome, v2DatabaseFile); sys.FileExists(v2DBFile) {
		result.found = true
		result.version = version2
	} else if v1DBFile := filepath.Join(defaultHome, v1DatabaseFile); sys.FileExists(v1DBFile) {
		result.found = true
	}

	return result
}
