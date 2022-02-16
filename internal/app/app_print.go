package app

import (
	"github.com/lemoony/snipkit/internal/ui"
)

func (a *appImpl) LookupAndCreatePrintableSnippet() (string, bool) {
	snippet := a.LookupSnippet()
	if snippet == nil {
		return "", false
	}

	parameters := snippet.GetParameters()
	if parameterValues, ok := a.tui.ShowParameterForm(parameters, ui.OkButtonPrint); ok {
		return snippet.Format(parameterValues, formatOptions(a.config.Script)), true
	}

	return "", false
}
