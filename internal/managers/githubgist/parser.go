package githubgist

import (
	"fmt"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/titleheader"
)

func parseSnippet(raw rawSnippet, cfg GistConfig) model.Snippet {
	result := model.Snippet{}
	result.UUID = raw.ID
	result.TagUUIDs = []string{} // TODO
	result.TitleFunc = func() string {
		if cfg.TitleHeaderEnabled {
			if title, ok := titleheader.ParseTitleFromHeader(string(raw.Content)); ok {
				return title
			}
		}

		switch cfg.NameMode {
		case SnippetNameModeDescription:
			return raw.Description
		case SnippetNameModeFilename:
			return raw.Filename
		case SnippetNameModeCombine:
			return fmt.Sprintf("%s - %s", raw.Description, raw.Filename)
		case SnippetNameModeCombinePreferDescription:
			if raw.FilesInGist == 1 {
				return raw.Description
			}
		case SnippetNameModeCombinePreferFilename:
			if raw.FilesInGist == 1 {
				return raw.Filename
			}
		}
		return fmt.Sprintf("%s - %s", raw.Description, raw.Filename)
	}
	result.ContentFunc = func() string {
		if cfg.HideTitleInPreview {
			return titleheader.PruneTitleHeader(string(raw.Content))
		}
		return string(raw.Content)
	}
	result.LanguageFunc = func() model.Language {
		return model.LanguageBash
	}
	return result
}
