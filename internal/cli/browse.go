package cli

import (
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/app"
	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/ui"
)

func LookupSnippet(v *viper.Viper) (*model.Snippet, error) {
	snipkit, err := app.NewApp(v)
	if snipkit == nil || err != nil {
		return nil, err
	}

	snippets, err := snipkit.GetAllSnippets()
	if err != nil {
		return nil, err
	}

	index, err := ui.ShowLookup(snippets)
	if index < 0 || err != nil {
		return nil, err
	}

	return &snippets[index], nil
}
