package app

import (
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
)

func (a *appImpl) LookupAndCreatePrintableSnippet() (bool, string) {
	if ok, snippet := a.LookupSnippet(); ok {
		parameters := snippet.GetParameters()
		if parameterValues, paramOk := a.tui.ShowParameterForm(parameters, nil, ui.OkButtonPrint); paramOk {
			return true, snippet.Format(parameterValues, formatOptions(a.config.Script))
		}
	}

	return false, ""
}

func (a *appImpl) LookupSnippetArgs() (bool, string, []model.ParameterValue) {
	if ok, snippet := a.LookupSnippet(); ok {
		parameters := snippet.GetParameters()
		if parameterValues, paramOk := a.tui.ShowParameterForm(parameters, nil, ui.OkButtonPrint); paramOk {
			return true, snippet.GetID(), matchParameterToValues(parameters, parameterValues)
		}
	}
	return false, "", nil
}

func (a *appImpl) FindSnippetAndPrint(id string, paramValues []model.ParameterValue) (bool, string) {
	if snippetFound, snippet := a.getSnippet(id); !snippetFound {
		panic(ErrSnippetIDNotFound)
	} else if paramOk, parameters := matchParameters(paramValues, snippet.GetParameters()); paramOk {
		return true, snippet.Format(parameters, formatOptions(a.config.Script))
	} else if selectedParams, formOk := a.tui.ShowParameterForm(snippet.GetParameters(), paramValues, ui.OkButtonExecute); formOk {
		return true, snippet.Format(selectedParams, formatOptions(a.config.Script))
	}
	return false, ""
}

func matchParameterToValues(parameters []model.Parameter, values []string) []model.ParameterValue {
	result := make([]model.ParameterValue, len(parameters))
	for i := range parameters {
		result[i] = model.ParameterValue{Key: parameters[i].Key, Value: values[i]}
	}
	return result
}
