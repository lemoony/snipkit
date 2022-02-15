package model

import (
	"fmt"
	"strings"
)

type Snippet struct {
	UUID     string
	TagUUIDs []string

	TitleFunc     func() string
	ContentFunc   func() string
	LanguageFunc  func() Language
	ParameterFunc func() []Parameter
	FormatFunc    func(content string, values []string) string
}

func (s *Snippet) GetTitle() string {
	return s.TitleFunc()
}

func (s *Snippet) GetContent() string {
	return s.ContentFunc()
}

func (s *Snippet) GetLanguage() Language {
	return s.LanguageFunc()
}

func (s Snippet) String() string {
	return fmt.Sprintf(
		"UUD: %s, Title: %s, Tags: [%s], Language: %d Content: %s",
		s.UUID,
		s.GetTitle(),
		strings.Join(s.TagUUIDs, ","),
		s.GetLanguage(),
		s.GetContent(),
	)
}

func (s Snippet) GetParameters() []Parameter {
	if s.ParameterFunc == nil {
		return nil
	}
	return s.ParameterFunc()
}

func (s Snippet) Format(values []string) string {
	if s.FormatFunc == nil {
		return ""
	}
	return s.FormatFunc(s.GetContent(), values)
}
