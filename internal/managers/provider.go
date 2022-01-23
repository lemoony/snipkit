package managers

import (
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Provider interface {
	CreateManager(system system.System, config Config) ([]Manager, error)
	ManagerDescriptions(config Config) []model.ManagerDescription
	AutoConfig(key model.ManagerKey, s *system.System) Config
}

type providerImpl struct{}

func NewBuilder() Provider {
	return providerImpl{}
}

func (p providerImpl) CreateManager(system system.System, config Config) ([]Manager, error) {
	var managers []Manager

	if config.SnippetsLab != nil {
		if manager, err := snippetslab.NewManager(
			snippetslab.WithSystem(&system),
			snippetslab.WithConfig(*config.SnippetsLab),
		); err != nil {
			return nil, err
		} else if manager != nil {
			managers = append(managers, manager)
		}
	}

	if config.PictarineSnip != nil {
		if manager, err := pictarinesnip.NewManager(
			pictarinesnip.WithSystem(&system),
			pictarinesnip.WithConfig(*config.PictarineSnip),
		); err != nil {
			return nil, err
		} else if manager != nil {
			managers = append(managers, manager)
		}
	}

	if config.FsLibrary != nil {
		if manager, err := fslibrary.NewManager(
			fslibrary.WithSystem(&system),
			fslibrary.WithConfig(*config.FsLibrary),
		); err != nil {
			return nil, err
		} else if manager != nil {
			managers = append(managers, manager)
		}
	}

	log.Info().Msgf("Number of enabled managers: %d", len(managers))

	return managers, nil
}

func (p providerImpl) ManagerDescriptions(config Config) []model.ManagerDescription {
	var infos []model.ManagerDescription
	if config.SnippetsLab == nil || !config.SnippetsLab.Enabled {
		infos = append(infos, snippetslab.Description(config.SnippetsLab))
	}
	if config.PictarineSnip == nil || !config.PictarineSnip.Enabled {
		infos = append(infos, pictarinesnip.Description(config.PictarineSnip))
	}
	if config.PictarineSnip == nil || !config.FsLibrary.Enabled {
		infos = append(infos, fslibrary.Description(config.FsLibrary))
	}
	return infos
}

func (p providerImpl) AutoConfig(key model.ManagerKey, s *system.System) Config {
	switch key {
	case snippetslab.Key:
		return Config{SnippetsLab: snippetslab.AutoDiscoveryConfig(s)}
	case pictarinesnip.Key:
		return Config{PictarineSnip: pictarinesnip.AutoDiscoveryConfig(s)}
	case fslibrary.Key:
		return Config{FsLibrary: fslibrary.AutoDiscoveryConfig(s)}
	}
	return Config{}
}