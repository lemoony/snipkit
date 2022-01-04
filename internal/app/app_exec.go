package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"emperror.dev/errors"

	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/parser"
)

func (a *appImpl) LookupAndExecuteSnippet() {
	snippet := a.LookupSnippet()
	if snippet == nil {
		return
	}

	parameters := parser.ParseParameters(snippet.Content)
	parameterValues := a.ui.ShowParameterForm(parameters)

	executeSnippet(*snippet, parameters, parameterValues)
}

func executeSnippet(snippet model.Snippet, params []model.Parameter, paramValues []string) {
	script := snippet.Content
	for i, p := range params {
		value := paramValues[i]
		if value == "" {
			value = p.DefaultValue
		}
		script = strings.ReplaceAll(script, fmt.Sprintf("${%s}", p.Key), value)
	}
	executeScript(script)
}

func executeScript(script string) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		panic(errors.New("no shell found"))
	}

	//nolint:gosec // since it would report G204 complaining about using a variable as input for exec.Command
	cmd, err := exec.Command(shell, "-c", script).Output()
	if err != nil {
		panic(err)
	}

	fmt.Print(string(cmd))
}
