package app

import (
	"github.com/lemoony/snippet-kit/internal/parser"
)

func (a *appImpl) LookupAndCreatePrintableSnippet() string {
	snippet := a.LookupSnippet()

	parameters := parser.ParseParameters(snippet.GetContent())
	parameterValues := a.ui.ShowParameterForm(parameters)

	return parser.CreateSnippet(snippet.GetContent(), parameters, parameterValues)
}
