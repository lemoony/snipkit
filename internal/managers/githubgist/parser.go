package githubgist

import (
	"fmt"

	"github.com/lemoony/snipkit/internal/model"
)

// TODO too many placeholders.
func parseSnippet(raw rawSnippet) model.Snippet {
	result := model.Snippet{}
	result.UUID = raw.ID
	result.TagUUIDs = []string{}
	result.TitleFunc = func() string {
		return fmt.Sprintf("%s - %s", raw.Description, raw.Filename)
	}
	result.ContentFunc = func() string {
		return string(raw.Content)
	}
	result.LanguageFunc = func() model.Language {
		return model.LanguageBash
	}
	return result
}
