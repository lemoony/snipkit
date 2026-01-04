package assistant

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"time"

	"github.com/lemoony/snipkit/internal/model"
)

// Pre-compiled regex patterns for script parsing.
var (
	scriptBlockRegex  = regexp.MustCompile("```(bash|sh)\\s+([\\s\\S]*?)```")
	filenameLineRegex = regexp.MustCompile(`(?m)^# Filename:\s*(\S+)\s*\n`)
	titleLineRegex    = regexp.MustCompile(`(?m)^# Snippet Title:\s*(.+)\s*\n`)
	commentCleanRegex = regexp.MustCompile(`(?m)^(\s*#\s*$\n\s*#\s*$\n?)`)
)

func PrepareSnippet(content []byte, parsed ParsedScript) model.Snippet {
	contentStr := string(content)
	return snippetImpl{
		id:      "",
		path:    "",
		tags:    []string{},
		content: contentStr,
		titleFunc: func() string {
			return parsed.Title
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
	// Find all matches of bash/sh code blocks
	matches := scriptBlockRegex.FindAllStringSubmatch(text, -1)

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

	// Step 1: Extract and remove the "# Filename:" line
	filenameMatch := filenameLineRegex.FindStringSubmatch(script)
	if len(filenameMatch) > 1 {
		filename = filenameMatch[1]
		removedLines++
	}
	script = filenameLineRegex.ReplaceAllString(script, "")

	// Step 2: Extract and remove the "# Snippet Title:" line
	titleMatch := titleLineRegex.FindStringSubmatch(script)
	if len(titleMatch) > 1 {
		title = titleMatch[1]
		removedLines++
	}
	script = titleLineRegex.ReplaceAllString(script, "")

	// Clean up empty comment lines if we removed metadata
	if removedLines > 0 {
		script = commentCleanRegex.ReplaceAllString(script, "")
	}

	return ParsedScript{
		Contents: script,
		Filename: filename,
		Title:    title,
	}
}
