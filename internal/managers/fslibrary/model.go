package fslibrary

import (
	"path/filepath"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/parser"
)

type snippetImpl struct {
	id          string
	path        string
	tags        []string
	titleFunc   func() string
	contentFunc func() string
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
	return s.contentFunc()
}

func (s snippetImpl) GetLanguage() model.Language {
	return languageForSuffix(filepath.Ext(s.path))
}

func (s snippetImpl) GetParameters() []model.Parameter {
	return parser.ParseParameters(s.GetContent())
}

func (s snippetImpl) Format(values []string, options model.SnippetFormatOptions) string {
	return parser.CreateSnippet(s.GetContent(), s.GetParameters(), values, options)
}
