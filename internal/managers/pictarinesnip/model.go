package pictarinesnip

import (
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/parser"
)

type snippetImpl struct {
	id       string
	tags     []string
	title    string
	language model.Language
	content  string
}

func (s snippetImpl) GetID() string {
	return s.id
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
	return parser.ParseParameters(s.content)
}

func (s snippetImpl) Format(values []string, options model.SnippetFormatOptions) string {
	return parser.CreateSnippet(s.GetContent(), s.GetParameters(), values, options)
}
