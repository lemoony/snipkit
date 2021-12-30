package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/parser"
	"github.com/lemoony/snippet-kit/internal/ui"
)

func LookupAndCreatePrintableSnippet(v *viper.Viper, term ui.Terminal) (string, error) {
	snippet, err := LookupSnippet(v, term)
	if snippet == nil || err != nil {
		return "", err
	}

	parameters := parser.ParseParameters(snippet.Content)
	parameterValues, err := term.ShowParameterForm(parameters)
	if err != nil {
		return "", err
	}

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
