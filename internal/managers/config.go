package managers

import (
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/githubgist"
	"github.com/lemoony/snipkit/internal/managers/pet"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
)

type Config struct {
	SnippetsLab   *snippetslab.Config   `yaml:"snippetsLab,omitempty" mapstructure:"snippetsLab"`
	PictarineSnip *pictarinesnip.Config `yaml:"pictarineSnip,omitempty" mapstructure:"pictarineSnip"`
	Pet           *pet.Config           `yaml:"pet,omitempty" mapstructure:"pet"`
	GithubGist    *githubgist.Config    `yaml:"githubGist,omitempty" mapstructure:"githubGist"`
	FsLibrary     *fslibrary.Config     `yaml:"fsLibrary,omitempty" mapstructure:"fsLibrary"`
}
