package app

import (
	"github.com/lemoony/snippet-kit/internal/model"
)

func (a *appImpl) LookupSnippet() *model.Snippet {
	snippets, err := a.getAllSnippets()
	if err != nil {
		// TODO: Handle no snippets
		panic(err)
	}

	if index := a.ui.ShowLookup(snippets); index < 0 {
		return nil
	} else {
		return &snippets[index]
	}
}
