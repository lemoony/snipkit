package ai

import (
	"regexp"

	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/ai/openai"
	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/titleheader"
)

const markdownScriptParts = 3

type Assistant interface {
	Query(string) string
}

type assistantImpl struct {
	system *system.System
	config Config
	cache  cache.Cache
}

func NewAssistant(system *system.System, config Config, cache cache.Cache) Assistant {
	return assistantImpl{system: system, config: config, cache: cache}
}

func (a assistantImpl) Query(prompt string) string {
	client, err := openai.NewClient(openai.WithCache(a.cache))
	if err != nil {
		panic(err)
	}

	response := client.Query(prompt)
	return extractBashScript(response)
}

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
