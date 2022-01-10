package testutil

import "github.com/lemoony/snippet-kit/internal/model"

func FixedString(title string) func() string {
	return func() string {
		return title
	}
}

func FixedLanguage(lang model.Language) func() model.Language {
	return func() model.Language {
		return lang
	}
}
