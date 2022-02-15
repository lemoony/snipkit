package pictarinesnip

import "github.com/lemoony/snipkit/internal/model"

type snippetImpl struct {
	id   string
	tags []string

	titleFunc     func() string
	contentFunc   func() string
	languageFunc  func() model.Language
	parameterFunc func() []model.Parameter
	formatFunc    func(content string, values []string) string
}

func (s snippetImpl) GetID() string {
	return "unused"
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
	return s.languageFunc()
}

func (s snippetImpl) GetParameters() []model.Parameter {
	return s.parameterFunc()
}

func (s snippetImpl) Format(values []string) string {
	return s.formatFunc(s.contentFunc(), values)
}

func (s snippetImpl) String() string {
	// TODO implement me
	panic("implement me")
}
