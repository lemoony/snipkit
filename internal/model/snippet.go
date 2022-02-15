package model

type Snippet interface {
	GetID() string
	GetTitle() string
	GetContent() string
	GetTags() []string
	GetLanguage() Language
	GetParameters() []Parameter
	Format([]string) string
	String() string
}
