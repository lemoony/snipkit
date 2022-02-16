package config

import (
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/ui"
)

type ParameterMode string

const (
	ParameterModeSet     = "SET"
	ParameterModeReplace = "Replace"
)

type VersionWrapper struct {
	Version string `yaml:"version" mapstructure:"version"`
	Config  Config `yaml:"config" mapstructure:"config"`
}

type Config struct {
	Style              ui.Config       `yaml:"style" mapstructure:"style"`
	Editor             string          `yaml:"editor" mapstructure:"editor" head_comment:"Your preferred editor to open the config file when typing 'snipkit config edit'." line_comment:"Defaults to a reasonable value for your operation system when empty."`
	DefaultRootCommand string          `yaml:"defaultRootCommand" mapstructure:"defaultRootCommand" head_comment:"The command which should run if you don't provide any subcommand." line_comment:"If not set, the help text will be shown."`
	Script             ScriptConfig    `yaml:"scripts" mapstructure:"Options regarding script handling"`
	Manager            managers.Config `yaml:"manager" mapstructure:"manager"`
}

type ScriptConfig struct {
	Shell          string        `yaml:"shell" mapstructure:"shell" head_comment:"The path to the shell to execute scripts with. If not set or empty, $SHELL will be used instead. Fallback is '/bin/bash'."`
	ParameterMode  ParameterMode `yaml:"parameterMode" mapstructure:"parameterMode" head_comment:"Defines how parameters are handled. Allowed values: SET (sets the parameter value as shell variable) and REPLACE (replaces all occurrences of the variable with the actual value)"`
	RemoveComments bool          `yaml:"removeComments" mapstructure:"removeComments" head_comment:"If set to true, any comments in your scripts will be removed upon executing or printing."`
}
