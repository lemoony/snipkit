package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/parser"
	"github.com/lemoony/snippet-kit/internal/ui"
)

func LookupAndExecuteSnippet() error {
	snippet, err := LookupSnippet()
	if err != nil {
		return err
	}

	parameters := parser.ParseParameters(snippet.Content)
	parameterValues := ui.ShowParameterForm(parameters)

	return executeSnippet(*snippet, parameters, parameterValues)
}

func executeSnippet(snippet model.Snippet, params []model.Parameter, paramValues []string) error {
	script := snippet.Content
	for i, p := range params {
		script = strings.ReplaceAll(script, fmt.Sprintf("${%s}", p.Key), paramValues[i])
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
