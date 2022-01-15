package providers

import (
	"github.com/lemoony/snipkit/internal/providers/fslibrary"
	"github.com/lemoony/snipkit/internal/providers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/providers/snippetslab"
)

type Config struct {
	SnippetsLab   snippetslab.Config   `yaml:"snippetsLab" mapstructure:"snippetsLab"`
	PictarineSnip pictarinesnip.Config `yaml:"pictarineSnip" mapstructure:"pictarineSnip"`
	FsLibrary     fslibrary.Config     `yaml:"fsLibrary" mapstructure:"fsLibrary"`
}
