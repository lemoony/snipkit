package providers

import (
	"github.com/lemoony/snippet-kit/internal/providers/filesystem"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
)

type Config struct {
	SnippetsLab snippetslab.Config `yaml:"snippetsLab" mapstructure:"snippetsLab"`
	FileSystem  filesystem.Config  `yaml:"fileSystem" mapstructure:"fileSystem"`
}
