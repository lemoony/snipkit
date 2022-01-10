package model

import (
	"fmt"
	"strings"
)

type Snippet struct {
	UUID     string
	title    string
	TagUUIDs []string

	TitleFunc    func() string
	ContentFunc  func() string
	LanguageFunc func() Language
}

func (s *Snippet) GetTitle() string {
	if s.TitleFunc != nil {
		return s.TitleFunc()
	}

	return s.title
}

func (s *Snippet) SetTitle(title string) {
	s.title = title
}

func (s *Snippet) GetContent() string {
	if s.ContentFunc != nil {
		return s.ContentFunc()
	}

	return s.title
}

func (s *Snippet) GetLanguage() Language {
	return s.LanguageFunc()
}

func (s Snippet) String() string {
	return fmt.Sprintf(
		"UUD: %s, Title: %s, Tags: [%s], Language: %d Content: %s",
		s.UUID,
		s.title,
		strings.Join(s.TagUUIDs, ","),
		s.GetLanguage(),
		s.GetContent(),
	)
}
