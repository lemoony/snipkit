package providers

import (
	"github.com/lemoony/snippet-kit/internal/providers/fslibrary"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
)

type Config struct {
	SnippetsLab snippetslab.Config `yaml:"snippetsLab" mapstructure:"snippetsLab"`
	FsLibrary   fslibrary.Config   `yaml:"fsLibrary" mapstructure:"fsLibrary"`
}
