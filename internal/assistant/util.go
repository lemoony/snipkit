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

func extractBashScript(text string) (string, string) {
	// Regex pattern to match bash script blocks in markdown
	pattern := "```(bash|sh)\\s+([\\s\\S]*?)```"
	re := regexp.MustCompile(pattern)

	// Find all matches of bash/sh code blocks
	matches := re.FindAllStringSubmatch(text, -1)

	var script string
	var filename string

	if len(matches) > 0 {
		// Extract the first matched bash/sh script block
		script = matches[0][2]
	} else {
		// If no markdown code block is found, assume the text is a bash script
		script = text
	}

	// Step 1: Remove the line starting with "# Filename:"
	// Use regular expressions to match and remove the entire line starting with "# Filename:"
	filenameLineRe := regexp.MustCompile(`(?m)^# Filename:\s*(\S+)\s*\n`)
	// Extract the filename if it exists
	filenameMatch := filenameLineRe.FindStringSubmatch(script)
	if len(filenameMatch) > 1 {
		filename = filenameMatch[1] // Extracted filename
	}

	// Remove the "# Filename:" line from the script
	scriptWithoutFilename := filenameLineRe.ReplaceAllString(script, "")

	return scriptWithoutFilename, filename
}
