package app

import (
	"fmt"
	"strings"

	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/parser"
)

func (a *appImpl) LookupAndCreatePrintableSnippet() string {
	snippet := a.LookupSnippet()

	parameters := parser.ParseParameters(snippet.Content)
	parameterValues := a.ui.ShowParameterForm(parameters)

	return createSnippetString(*snippet, parameters, parameterValues)
}

func createSnippetString(snippet model.Snippet, params []model.Parameter, paramValues []string) string {
	script := snippet.Content
	for i, p := range params {
		value := paramValues[i]
		if value == "" {
			value = p.DefaultValue
		}
		script = strings.ReplaceAll(script, fmt.Sprintf("${%s}", p.Key), value)
	}
	return script
}
