package managers

import (
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/githubgist"
	"github.com/lemoony/snipkit/internal/managers/masscode"
	"github.com/lemoony/snipkit/internal/managers/pet"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Provider interface {
	CreateManager(system system.System, config Config) []Manager
	ManagerDescriptions(config Config) []model.ManagerDescription
	AutoConfig(key model.ManagerKey, s *system.System) Config
}

type providerImpl struct {
	cache cache.Cache
}

func NewBuilder(cache cache.Cache) Provider {
	return providerImpl{cache: cache}
}

func (p providerImpl) CreateManager(system system.System, config Config) []Manager {
	var managers []Manager

	if manager := createSnippetsLab(system, config); manager != nil {
		managers = append(managers, manager)
	}
	if manager := createPictarineSnip(system, config); manager != nil {
		managers = append(managers, manager)
	}
	if manager := createPetConfig(system, config); manager != nil {
		managers = append(managers, manager)
	}
	if manager := createMassCodeConfig(system, config); manager != nil {
		managers = append(managers, manager)
	}
	if manager := createGitHubGist(system, config, p.cache); manager != nil {
		managers = append(managers, manager)
	}
	if manager := createFSLibrary(system, config); manager != nil {
		managers = append(managers, manager)
	}

	log.Info().Msgf("Number of enabled managers: %d", len(managers))

	return managers
}

func (p providerImpl) ManagerDescriptions(config Config) []model.ManagerDescription {
	var infos []model.ManagerDescription
	if config.SnippetsLab == nil || !config.SnippetsLab.Enabled {
		infos = append(infos, snippetslab.Description(config.SnippetsLab))
	}
	if config.PictarineSnip == nil || !config.PictarineSnip.Enabled {
		infos = append(infos, pictarinesnip.Description(config.PictarineSnip))
	}
	if config.Pet == nil || !config.Pet.Enabled {
		infos = append(infos, pet.Description(config.Pet))
	}
	if config.MassCode == nil || !config.MassCode.Enabled {
		infos = append(infos, masscode.Description(config.MassCode))
	}
	if config.GithubGist == nil || !config.GithubGist.Enabled {
		infos = append(infos, githubgist.Description(config.GithubGist))
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
	case pet.Key:
		return Config{Pet: pet.AutoDiscoveryConfig(s)}
	case masscode.Key:
		return Config{MassCode: masscode.AutoDiscoveryConfig(s)}
	case githubgist.Key:
		return Config{GithubGist: githubgist.AutoDiscoveryConfig()}
	case fslibrary.Key:
		return Config{FsLibrary: fslibrary.AutoDiscoveryConfig(s)}
	}
	return Config{}
}

func createSnippetsLab(system system.System, config Config) Manager {
	if config.SnippetsLab == nil || !config.SnippetsLab.Enabled {
		return nil
	}
	manager, err := snippetslab.NewManager(
		snippetslab.WithSystem(&system),
		snippetslab.WithConfig(*config.SnippetsLab),
	)
	if err != nil {
		panic(err)
	}
	return manager
}

func createPictarineSnip(system system.System, config Config) Manager {
	if config.PictarineSnip == nil || !config.PictarineSnip.Enabled {
		return nil
	}
	manager, err := pictarinesnip.NewManager(
		pictarinesnip.WithSystem(&system),
		pictarinesnip.WithConfig(*config.PictarineSnip),
	)
	if err != nil {
		panic(err)
	}
	return manager
}

func createPetConfig(system system.System, config Config) Manager {
	if config.Pet == nil || !config.Pet.Enabled {
		return nil
	}
	manager, err := pet.NewManager(pet.WithSystem(&system), pet.WithConfig(*config.Pet))
	if err != nil {
		panic(err)
	}
	return manager
}

func createMassCodeConfig(system system.System, config Config) Manager {
	if config.Pet == nil || !config.Pet.Enabled {
		return nil
	}
	manager, err := masscode.NewManager(masscode.WithSystem(&system), masscode.WithConfig(*config.MassCode))
	if err != nil {
		panic(err)
	}
	return manager
}

func createGitHubGist(system system.System, config Config, cache cache.Cache) Manager {
	if config.GithubGist == nil || !config.GithubGist.Enabled {
		return nil
	}
	manager, err := githubgist.NewManager(
		githubgist.WithSystem(&system),
		githubgist.WithConfig(*config.GithubGist),
		githubgist.WithCache(cache),
	)
	if err != nil {
		panic(err)
	}
	return manager
}

func createFSLibrary(system system.System, config Config) Manager {
	if config.FsLibrary == nil || !config.FsLibrary.Enabled {
		return nil
	}
	manager, err := fslibrary.NewManager(fslibrary.WithSystem(&system), fslibrary.WithConfig(*config.FsLibrary))
	if err != nil {
		panic(err)
	}
	return manager
}
