package app

import (
	"github.com/lemoony/snipkit/internal/parser"
	"github.com/lemoony/snipkit/internal/ui"
)

func (a *appImpl) LookupAndCreatePrintableSnippet() (string, bool) {
	snippet := a.LookupSnippet()
	if snippet == nil {
		return "", false
	}

	parameters := parser.ParseParameters(snippet.GetContent())
	if parameterValues, ok := a.ui.ShowParameterForm(parameters, ui.OkButtonPrint); ok {
		return parser.CreateSnippet(snippet.GetContent(), parameters, parameterValues), true
	}

	return "", false
}
