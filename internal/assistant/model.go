package assistant

import (
	"path/filepath"

	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/parser"
)

type snippetImpl struct {
	id        string
	path      string
	content   string
	tags      []string
	titleFunc func() string
}

func (s snippetImpl) GetID() string {
	return s.id
}

func (s snippetImpl) GetTitle() string {
	return s.titleFunc()
}

func (s snippetImpl) GetTags() []string {
	return s.tags
}

func (s snippetImpl) GetContent() string {
	return s.content
}

func (s snippetImpl) GetLanguage() model.Language {
	return fslibrary.LanguageForSuffix(filepath.Ext(s.path))
}

func (s snippetImpl) GetParameters() []model.Parameter {
	return parser.ParseParameters(s.GetContent())
}

func (s snippetImpl) Format(values []string, options model.SnippetFormatOptions) string {
	return parser.CreateSnippet(s.GetContent(), s.GetParameters(), values, options)
}
