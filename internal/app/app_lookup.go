package app

import (
	"github.com/lemoony/snippet-kit/internal/model"
)

func (a *appImpl) LookupSnippet() *model.Snippet {
	snippets, err := a.getAllSnippets()
	if err != nil {
		panic(err)
	}

	index := a.ui.ShowLookup(snippets)

	return &snippets[index]
}
