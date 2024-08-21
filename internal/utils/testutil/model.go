package testutil

import (
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/parser"
)

type TestSnippet struct {
	ID   string
	Tags []string

	Title    string
	Content  string
	Language model.Language
}

var DummySnippet = TestSnippet{ID: "uuid-x", Title: "title-2", Language: model.LanguageBash, Tags: []string{}, Content: "testSnippetContent"}

func (t TestSnippet) GetID() string {
	return t.ID
}

func (t TestSnippet) GetTitle() string {
	return t.Title
}

func (t TestSnippet) GetContent() string {
	return t.Content
}

func (t TestSnippet) GetTags() []string {
	return t.Tags
}

func (t TestSnippet) GetLanguage() model.Language {
	return t.Language
}

func (t TestSnippet) GetParameters() []model.Parameter {
	return parser.ParseParameters(t.Content)
}

func (t TestSnippet) Format(values []string, options model.SnippetFormatOptions) string {
	return parser.CreateSnippet(t.Content, t.GetParameters(), values, options)
}

func (t TestSnippet) String() string {
	panic("implement me")
}
