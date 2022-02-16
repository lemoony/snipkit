package model

type SnippetParamMode int

const (
	SnippetParamModeSet     = 0
	SnippetParamModeReplace = 1
)

type SnippetFormatOptions struct {
	RemoveComments bool
	ParamMode      SnippetParamMode
}

type Snippet interface {
	GetID() string
	GetTitle() string
	GetContent() string
	GetTags() []string
	GetLanguage() Language
	GetParameters() []Parameter
	Format([]string, SnippetFormatOptions) string
}
