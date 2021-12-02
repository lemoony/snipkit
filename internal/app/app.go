package app

import (
	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/providers"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
	"github.com/lemoony/snippet-kit/internal/utils"
)

type App struct {
	Providers []providers.Provider
}

func NewApp() (*App, error) {
	system, err := utils.NewSystem()
	if err != nil {
		return nil, err
	}

	snippetsLab, err := snippetslab.NewProvider(
		snippetslab.WithSystem(&system),
		// snippetslab.WithTagsFilter([]string{"snipkit", "footag"}),
	)
	if err != nil {
		return nil, err
	}

	allProviders := []providers.Provider{
		snippetsLab,
	}

	return &App{
		Providers: allProviders,
	}, nil
}

func (a *App) GetAllSnippets() ([]model.Snippet, error) {
	var result []model.Snippet
	for _, provider := range a.Providers {
		if snippets, err := provider.GetSnippets(); err != nil {
			return nil, err
		} else {
			result = append(result, snippets...)
		}
	}
	return result, nil
}
