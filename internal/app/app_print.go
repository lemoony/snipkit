package app

import (
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
)

func (a *appImpl) LookupAndCreatePrintableSnippet() (string, bool) {
	snippet := a.LookupSnippet()
	if snippet == nil {
		return "", false
	}

	parameters := snippet.GetParameters()
	if parameterValues, ok := a.tui.ShowParameterForm(parameters, nil, ui.OkButtonPrint); ok {
		return snippet.Format(parameterValues, formatOptions(a.config.Script)), true
	}

	return "", false
}

func (a *appImpl) FindSnippetAndPrint(id string, paramValues []model.ParameterValue) (string, bool) {
	if snippetFound, snippet := a.getSnippet(id); !snippetFound {
		panic(ErrSnippetIDNotFound)
	} else if paramOk, parameters := matchParameters(paramValues, snippet.GetParameters()); paramOk {
		return snippet.Format(parameters, formatOptions(a.config.Script)), true
	} else if selectedParams, formOk := a.tui.ShowParameterForm(snippet.GetParameters(), paramValues, ui.OkButtonExecute); formOk {
		return snippet.Format(selectedParams, formatOptions(a.config.Script)), true
	}
	return "", false
}
