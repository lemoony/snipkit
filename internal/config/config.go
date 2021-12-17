package config

import (
	"github.com/lemoony/snippet-kit/internal/providers/filesystem"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
	"github.com/lemoony/snippet-kit/internal/ui"
)

type VersionWrapper struct {
	Version string `yaml:"version" mapstructure:"version"`
	Config  Config `yaml:"config" mapstructure:"config"`
}

type Config struct {
	Style     ui.Config `yaml:"style" mapstructure:"style"`
	Editor    string    `yaml:"editor" mapstructure:"editor" head_comment:"Your preferred editor to open the config file when typing 'snipkit config edit'." line_comment:"Defaults to a reasonable value for your operation system when empty."`
	Providers Providers `yaml:"provider" mapstructure:"provider"`
}

type Providers struct {
	SnippetsLab snippetslab.Config `yaml:"snippetsLab" mapstructure:"snippetsLab"`
	FileSystem  filesystem.Config  `yaml:"fileSystem" mapstructure:"fileSystem"`
}
