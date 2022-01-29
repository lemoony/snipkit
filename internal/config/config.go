package config

import (
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/ui"
)

type VersionWrapper struct {
	Version string `yaml:"version" mapstructure:"version"`
	Config  Config `yaml:"config" mapstructure:"config"`
}

type Config struct {
	Style              ui.Config       `yaml:"style" mapstructure:"style"`
	Editor             string          `yaml:"editor" mapstructure:"editor" head_comment:"Your preferred editor to open the config file when typing 'snipkit config edit'." line_comment:"Defaults to a reasonable value for your operation system when empty."`
	DefaultRootCommand string          `yaml:"defaultRootCommand" mapstructure:"defaultRootCommand" head_comment:"The command which should run if you don't provide any subcommand." line_comment:"If not set, the help text will be shown."`
	Manager            managers.Config `yaml:"manager" mapstructure:"manager"`
}
