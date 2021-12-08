package cli

import (
	"fmt"
	"strings"

	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/parser"
	"github.com/lemoony/snippet-kit/internal/ui"
)

func LookupAndCreatePrintableSnippet() (string, error) {
	snippet, err := LookupSnippet()
	if err != nil {
		return "", err
	}

	parameters := parser.ParseParameters(snippet.Content)
	parameterValues := ui.ShowParameterForm(parameters)

	return createSnippetString(*snippet, parameters, parameterValues)
}

func createSnippetString(snippet model.Snippet, params []model.Parameter, paramValues []string) (string, error) {
	script := snippet.Content
	for i, p := range params {
		value := paramValues[i]
		if value == "" {
			value = p.DefaultValue
		}
		script = strings.ReplaceAll(script, fmt.Sprintf("${%s}", p.Key), value)
	}
	return script, nil
}
