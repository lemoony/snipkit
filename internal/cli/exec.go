package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/parser"
	"github.com/lemoony/snippet-kit/internal/ui"
)

func LookupAndExecuteSnippet(v *viper.Viper) error {
	snippet, err := LookupSnippet(v)
	if snippet == nil || err != nil {
		return err
	}

	parameters := parser.ParseParameters(snippet.Content)
	parameterValues, err := ui.ShowParameterForm(parameters)
	if err != nil {
		return err
	}

	return executeSnippet(*snippet, parameters, parameterValues)
}

func executeSnippet(snippet model.Snippet, params []model.Parameter, paramValues []string) error {
	script := snippet.Content
	for i, p := range params {
		value := paramValues[i]
		if value == "" {
			value = p.DefaultValue
		}
		script = strings.ReplaceAll(script, fmt.Sprintf("${%s}", p.Key), value)
	}
	return executeScript(script)
}

func executeScript(script string) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return errors.New("no shell found")
	}

	//nolint:gosec // since it would report G204 complaining about using a variable as input for exec.Command
	cmd, err := exec.Command(shell, "-c", script).Output()
	if err != nil {
		return err
	}

	fmt.Print(string(cmd))
	return nil
}
