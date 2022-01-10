package app

import (
	"github.com/lemoony/snippet-kit/internal/model"
)

func (a *appImpl) LookupSnippet() *model.Snippet {
	snippets := a.getAllSnippets()
	if len(snippets) == 0 {
		panic(ErrNoSnippetsAvailable)
	}

	if index := a.ui.ShowLookup(snippets); index < 0 {
		return nil
	} else {
		return &snippets[index]
	}
}
