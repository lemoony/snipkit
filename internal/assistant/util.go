package assistant

import (
	"regexp"

	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/titleheader"
)

func PrepareSnippet(content string) model.Snippet {
	return snippetImpl{
		id:      "",
		path:    "",
		tags:    []string{},
		content: content,
		titleFunc: func() string {
			if title, ok := titleheader.ParseTitleFromHeader(content); ok {
				return title
			}
			return ""
		},
	}
}

func extractBashScript(text string) string {
	// Regex pattern to match bash script blocks in markdown
	pattern := "```(bash|sh)\\s+([\\s\\S]*?)```"
	re := regexp.MustCompile(pattern)

	// Find all matches of bash/sh code blocks
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) >= markdownScriptParts {
			return match[2]
		}
	}

	panic(errors.Errorf("Invalid response from AI provider: %s", text))
}
