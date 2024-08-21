package app

import (
	"github.com/lemoony/snipkit/internal/model"
)

func (a *appImpl) LookupSnippet() (bool, model.Snippet) {
	snippets := a.getAllSnippets()
	if len(snippets) == 0 {
		panic(ErrNoSnippetsAvailable)
	}

	if index := a.tui.ShowLookup(snippets, a.config.FuzzySearch); index < 0 {
		return false, nil
	} else {
		return true, snippets[index]
	}
}
