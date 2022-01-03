package cli

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/app"
	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/ui"
)

func LookupSnippet(v *viper.Viper, term ui.Terminal) (*model.Snippet, error) {
	snipkit, err := app.NewApp(v)
	if snipkit == nil || err != nil {
		return nil, err
	}

	snippets, err := snipkit.GetAllSnippets()
	if err != nil {
		return nil, err
	}

	index := term.ShowLookup(snippets)
	if index < 0 {
		panic(fmt.Sprintf("invalid index: %d", index))
	}

	return &snippets[index], nil
}
