package managers

import (
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
)

type Config struct {
	SnippetsLab   *snippetslab.Config   `yaml:"snippetsLab,omitempty" mapstructure:"snippetsLab"`
	PictarineSnip *pictarinesnip.Config `yaml:"pictarineSnip,omitempty" mapstructure:"pictarineSnip"`
	FsLibrary     *fslibrary.Config     `yaml:"fsLibrary,omitempty" mapstructure:"fsLibrary"`
}
