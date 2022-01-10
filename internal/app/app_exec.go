package app

import (
	"os"
	"os/exec"

	"github.com/lemoony/snippet-kit/internal/parser"
	"github.com/lemoony/snippet-kit/internal/ui"
	"github.com/lemoony/snippet-kit/internal/utils/stringutil"
)

func (a *appImpl) LookupAndExecuteSnippet() {
	snippet := a.LookupSnippet()
	if snippet == nil {
		return
	}

	parameters := parser.ParseParameters(snippet.GetContent())
	if parameterValues, ok := a.ui.ShowParameterForm(parameters); ok {
		script := parser.CreateSnippet(snippet.GetContent(), parameters, parameterValues)
		executeScript(script, a.ui)
	}
}

func executeScript(script string, term ui.Terminal) {
	shell := stringutil.StringOrDefault(os.Getenv("SHELL"), "/bin/bash")

	//nolint:gosec // since it would report G204 complaining about using a variable as input for exec.Command
	cmd, err := exec.Command(shell, "-c", script).CombinedOutput()
	if err != nil {
		panic(err)
	}

	term.PrintMessage(string(cmd))
}
