package managers

import (
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type ManagerKey int

const (
	SnippetsLabManagerKey   = ManagerKey(iota)
	PictarineSnipManagerKey = ManagerKey(iota)
	FsLibraryManagerKey     = ManagerKey(iota)
)

type Provider interface {
	CreateManager(system system.System, config Config) ([]Manager, error)
	ManagerDescriptions(config Config) []model.ManagerDescription
}

type providerImpl struct{}

func NewBuilder() Provider {
	return providerImpl{}
}

func (b providerImpl) CreateManager(system system.System, config Config) ([]Manager, error) {
	var managers []Manager

	if manager, err := snippetslab.NewManager(
		snippetslab.WithSystem(&system),
		snippetslab.WithConfig(config.SnippetsLab),
	); err != nil {
		return nil, err
	} else if manager != nil {
		managers = append(managers, manager)
	}

	if manager, err := fslibrary.NewManager(
		fslibrary.WithSystem(&system),
		fslibrary.WithConfig(config.FsLibrary),
	); err != nil {
		return nil, err
	} else if manager != nil {
		managers = append(managers, manager)
	}

	if manager, err := pictarinesnip.NewManager(
		pictarinesnip.WithSystem(&system),
		pictarinesnip.WithConfig(config.PictarineSnip),
	); err != nil {
		return nil, err
	} else if manager != nil {
		managers = append(managers, manager)
	}

	log.Info().Msgf("Number of enabled managers: %d", len(managers))

	return managers, nil
}

func (b providerImpl) ManagerDescriptions(config Config) []model.ManagerDescription {
	var infos []model.ManagerDescription
	infos = append(infos, snippetslab.Description(config.SnippetsLab))
	infos = append(infos, pictarinesnip.Description(config.PictarineSnip))
	infos = append(infos, fslibrary.Description(config.FsLibrary))
	return infos
}
