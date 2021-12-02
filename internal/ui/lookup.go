package ui

import (
	"fmt"

	"github.com/ktr0731/go-fuzzyfinder"

	"github.com/lemoony/snippet-kit/internal/model"
)

func ShowLookup(snippets []model.Snippet) (int, error) {
	return fuzzyfinder.Find(
		snippets,
		func(i int) string {
			return snippets[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Title: %s\n\nContent: %s",
				snippets[i].Title,
				snippets[i].Content,
			)
		}),
	)
}
