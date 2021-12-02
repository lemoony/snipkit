package cli

import (
	"github.com/lemoony/snippet-kit/internal/app"
	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/ui"
)

func LookupSnippet() (*model.Snippet, error) {
	snipkit, err := app.NewApp()
	if err != nil {
		return nil, err
	}

	snippets, err := snipkit.GetAllSnippets()
	if err != nil {
		return nil, err
	}

	index, err := ui.ShowLookup(snippets)
	if err != nil {
		return nil, err
	}

	return &snippets[index], nil
}
