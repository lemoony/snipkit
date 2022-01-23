package app

import (
	"github.com/lemoony/snipkit/internal/model"
)

func (a *appImpl) LookupSnippet() *model.Snippet {
	snippets := a.getAllSnippets()
	if len(snippets) == 0 {
		panic(ErrNoSnippetsAvailable)
	}

	if index := a.tui.ShowLookup(snippets); index < 0 {
		return nil
	} else {
		return &snippets[index]
	}
}
