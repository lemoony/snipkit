package model

type Language int

const (
	LanguageUnknown  = Language(0)
	LanguageBash     = Language(1)
	LanguageYAML     = Language(2)
	LanguageMarkdown = Language(3)
	LanguageText     = Language(4)
	LanguageTOML     = Language(5)
)
