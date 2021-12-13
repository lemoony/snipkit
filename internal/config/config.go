package config

import (
	"github.com/lemoony/snippet-kit/internal/providers/filesystem"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
)

type Configer interface {
	getTheme() string
	getEditor() string
}

type versionWrapper struct {
	Version string `yaml:"version"`
	Config  Config `yaml:"config"`
}

type Config struct {
	Style struct {
		Theme string `yaml:"theme" head_comment:"The theme defines the terminal colors used by Snipkit.\nAvailable themes:dracula,solarized,light,dark."`
	} `yaml:"style"`
	Editor    string `yaml:"editor" head_comment:"Your preferred editor to open the config file when typing 'snipkit config edit'." line_comment:"Defaults to a reasonable value for your operation system when empty."`
	Providers struct {
		SnippetsLab snippetslab.Config `yaml:"snippetsLab"`
		FileSystem  filesystem.Config  `yaml:"fileSystem"`
	} `yaml:"provider"`
}
