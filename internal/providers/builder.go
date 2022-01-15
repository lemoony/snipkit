package providers

import (
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/providers/fslibrary"
	"github.com/lemoony/snipkit/internal/providers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/providers/snippetslab"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Builder interface {
	BuildProvider(system system.System, config Config) ([]Provider, error)
}

type builderImpl struct{}

func NewBuilder() Builder {
	return builderImpl{}
}

func (b builderImpl) BuildProvider(system system.System, config Config) ([]Provider, error) {
	var providers []Provider

	if provider, err := snippetslab.NewProvider(
		snippetslab.WithSystem(&system),
		snippetslab.WithConfig(config.SnippetsLab),
	); err != nil {
		return nil, err
	} else if provider != nil {
		providers = append(providers, provider)
	}

	if provider, err := fslibrary.NewProvider(
		fslibrary.WithSystem(&system),
		fslibrary.WithConfig(config.FsLibrary),
	); err != nil {
		return nil, err
	} else if provider != nil {
		providers = append(providers, provider)
	}

	if provider, err := pictarinesnip.NewProvider(
		pictarinesnip.WithSystem(&system),
		pictarinesnip.WithConfig(config.PictarineSnip),
	); err != nil {
		return nil, err
	} else if provider != nil {
		providers = append(providers, provider)
	}

	log.Info().Msgf("Number of enabled providers: %d", len(providers))

	return providers, nil
}
