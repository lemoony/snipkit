package app

import (
	"os"
	"os/exec"

	"github.com/lemoony/snippet-kit/internal/parser"
	"github.com/lemoony/snippet-kit/internal/ui"
)

func (a *appImpl) LookupAndExecuteSnippet() {
	snippet := a.LookupSnippet()
	if snippet == nil {
		return
	}

	parameters := parser.ParseParameters(snippet.Content)
	parameterValues := a.ui.ShowParameterForm(parameters)

	script := createSnippetString(*snippet, parameters, parameterValues)

	executeScript(script, a.ui)
}

func executeScript(script string, term ui.Terminal) {
	shell := os.Getenv("SHELL")

	//nolint:gosec // since it would report G204 complaining about using a variable as input for exec.Command
	cmd, err := exec.Command(shell, "-c", script).CombinedOutput()
	if err != nil {
		panic(err)
	}

	term.PrintMessage(string(cmd))
}
