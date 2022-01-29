package app

import (
	"os"
	"os/exec"

	"emperror.dev/errors"
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/parser"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

func (a *appImpl) LookupAndExecuteSnippet() {
	snippet := a.LookupSnippet()
	if snippet == nil {
		return
	}

	parameters := parser.ParseParameters(snippet.GetContent())
	if parameterValues, ok := a.tui.ShowParameterForm(parameters, ui.OkButtonExecute); ok {
		script := parser.CreateSnippet(snippet.GetContent(), parameters, parameterValues)
		executeScript(script)
	}
}

func executeScript(script string) {
	shell := stringutil.StringOrDefault(os.Getenv("SHELL"), "/bin/bash")

	//nolint:gosec // since it would report G204 complaining about using a variable as input for exec.Command
	cmd := exec.Command(shell, "-c", script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		panic(errors.Wrapf(errors.WithStack(err), "failed to run command"))
	}

	if err := cmd.Wait(); err != nil {
		log.Info().Err(err)
	}
}
