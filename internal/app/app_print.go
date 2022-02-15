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

	parameters := snippet.GetParameters()
	if parameters == nil {
		parameters = parser.ParseParameters(snippet.GetContent())
	}
	if parameterValues, ok := a.tui.ShowParameterForm(parameters, ui.OkButtonPrint); ok {
		script := snippet.Format(parameterValues)
		if script == "" {
			script = parser.CreateSnippet(snippet.GetContent(), parameters, parameterValues)
		}
		return script, true
	}

	return "", false
}
