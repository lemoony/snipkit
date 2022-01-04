package app

import (
	"fmt"

	"github.com/lemoony/snippet-kit/internal/model"
)

func (a *appImpl) LookupSnippet() *model.Snippet {
	snippets, err := a.getAllSnippets()
	if err != nil {
		panic(err)
	}

	index := a.ui.ShowLookup(snippets)
	if index < 0 {
		panic(fmt.Sprintf("invalid index: %d", index))
	}

	return &snippets[index]
}
