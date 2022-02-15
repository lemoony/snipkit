package githubgist

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/titleheader"
)

var tagRegex = regexp.MustCompile(`#\S+`)

var languageMapping = map[string]model.Language{
	"Shell":    model.LanguageBash,
	"Markdown": model.LanguageMarkdown,
	"TOML":     model.LanguageTOML,
	"YAML":     model.LanguageYAML,
}

func parseSnippet(raw rawSnippet, cfg GistConfig) model.Snippet {
	result := snippetImpl{
		id:   raw.ID,
		tags: parseTags(raw.Description),
		titleFunc: func() string {
			return parseTitle(raw, cfg.NameMode, cfg.TitleHeaderEnabled)
		},
		contentFunc: func() string {
			return formatContent(string(raw.Content), cfg.HideTitleInPreview)
		},
		languageFunc: func() model.Language {
			return mapLanguage(raw.Language)
		},
	}
	return &result
}

func parseTitle(raw rawSnippet, nameMode SnippetNameMode, titleHeaderEnabled bool) string {
	if titleHeaderEnabled {
		if title, ok := titleheader.ParseTitleFromHeader(string(raw.Content)); ok {
			return title
		}
	}

	switch nameMode {
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
	}
	return fmt.Sprintf("%s - %s", raw.Description, raw.Filename)
}

func parseTags(text string) []string {
	tags := tagRegex.FindAllString(text, -1)
	for i := range tags {
		tags[i] = tags[i][1:]
	}
	if len(tags) == 0 {
		return []string{}
	}
	return tags
}

func pruneTags(text string) string {
	return strings.TrimSpace(tagRegex.ReplaceAllString(text, ""))
}

func formatContent(text string, hideTitleHeader bool) string {
	if hideTitleHeader {
		return titleheader.PruneTitleHeader(text)
	}
	return text
}

func mapLanguage(val string) model.Language {
	if lang, ok := languageMapping[val]; ok {
		return lang
	}
	return model.LanguageText
}
