package assistant

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"time"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/titleheader"
)

func PrepareSnippet(content []byte) model.Snippet {
	contentStr := string(content)
	return snippetImpl{
		id:      "",
		path:    "",
		tags:    []string{},
		content: contentStr,
		titleFunc: func() string {
			if title, ok := titleheader.ParseTitleFromHeader(contentStr); ok {
				return title
			}
			return ""
		},
	}
}

func RandomScriptFilename() string {
	const randomNumberLength = 5
	timestamp := time.Now().Format("20060102_150405")
	randomBytes := make([]byte, randomNumberLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	randomString := fmt.Sprintf("%x", randomBytes)
	return fmt.Sprintf("%s_%s.sh", timestamp, randomString)
}

type ParsedScript struct {
	Contents string
	Filename string
	Title    string
}

func parseScript(text string) ParsedScript {
	// Regex pattern to match bash Contents blocks in markdown
	pattern := "```(bash|sh)\\s+([\\s\\S]*?)```"
	re := regexp.MustCompile(pattern)

	// Find all matches of bash/sh code blocks
	matches := re.FindAllStringSubmatch(text, -1)

	var script string
	var filename string
	var title string

	if len(matches) > 0 {
		// Extract the first matched bash/sh Contents block
		script = matches[0][2]
	} else {
		// If no markdown code block is found, assume the text is a bash Contents
		script = text
	}

	removedLines := 0
	// Step 1: Remove the line starting with "# Filename:"
	filenameLineRe := regexp.MustCompile(`(?m)^# Filename:\s*(\S+)\s*\n`)
	// Extract the Filename if it exists
	filenameMatch := filenameLineRe.FindStringSubmatch(script)
	if len(filenameMatch) > 1 {
		filename = filenameMatch[1] // Extracted Filename
		removedLines++
	}
	// Remove the "# Filename:" line from the Contents
	script = filenameLineRe.ReplaceAllString(script, "")

	// Step 2: Remove the line starting with "# Snippet Title:"
	titleLineRe := regexp.MustCompile(`(?m)^# Snippet Title:\s*(.+)\s*\n`)
	// Extract the Snippet Title if it exists
	titleMatch := titleLineRe.FindStringSubmatch(script)
	if len(titleMatch) > 1 {
		title = titleMatch[1] // Extracted Title
		removedLines++
	}
	// Remove the "# Snippet Title:" line from the Contents
	script = titleLineRe.ReplaceAllString(script, "")

	if removedLines > 0 {
		commentRe := regexp.MustCompile(`(?m)^(\s*#\s*$\n\s*#\s*$\n?)`)
		script = commentRe.ReplaceAllString(script, "")
	}

	return ParsedScript{
		Contents: script,
		Filename: filename,
		Title:    title,
	}
}
