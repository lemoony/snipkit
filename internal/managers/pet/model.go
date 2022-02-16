package pet

import (
	"github.com/lemoony/snipkit/internal/model"
)

type snippetImpl struct {
	id       string
	tags     []string
	title    string
	content  string
	language model.Language
}

func (s snippetImpl) GetID() string {
	return "unused"
}

func (s snippetImpl) GetTitle() string {
	return s.title
}

func (s snippetImpl) GetTags() []string {
	return s.tags
}

func (s snippetImpl) GetContent() string {
	return s.content
}

func (s snippetImpl) GetLanguage() model.Language {
	return s.language
}

func (s snippetImpl) GetParameters() []model.Parameter {
	return parseParameters(s.content)
}

func (s snippetImpl) Format(values []string, _ model.SnippetFormatOptions) string {
	return formatContent(s.content, values)
}
