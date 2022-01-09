package providers

import (
	"github.com/lemoony/snippet-kit/internal/providers/fslibrary"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
	"github.com/lemoony/snippet-kit/internal/utils"
)

type Builder interface {
	BuildProvider(system utils.System, config Config) ([]Provider, error)
}

type builderImpl struct{}

func NewBuilder() Builder {
	return builderImpl{}
}

func (b builderImpl) BuildProvider(system utils.System, config Config) ([]Provider, error) {
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

	return providers, nil
}
