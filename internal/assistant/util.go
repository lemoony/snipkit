package assistant

import (
	"regexp"

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

	if len(matches) > 0 {
		// Return the first matched code block
		return matches[0][2]
	}

	// If no markdown code block is found, assume the text is a bash script
	return text
}
