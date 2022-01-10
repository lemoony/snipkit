package app

import (
	"github.com/lemoony/snippet-kit/internal/parser"
)

func (a *appImpl) LookupAndCreatePrintableSnippet() (string, bool) {
	snippet := a.LookupSnippet()
	if snippet == nil {
		return "", false
	}

	parameters := parser.ParseParameters(snippet.GetContent())
	if parameterValues, ok := a.ui.ShowParameterForm(parameters); ok {
		return parser.CreateSnippet(snippet.GetContent(), parameters, parameterValues), true
	}

	return "", false
}
