package app

import (
	"os"
	"os/exec"

	"emperror.dev/errors"
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

const fallbackShell = "/bin/bash"

func (a *appImpl) LookupAndExecuteSnippet() {
	snippet := a.LookupSnippet()
	if snippet == nil {
		return
	}

	parameters := snippet.GetParameters()
	if parameterValues, ok := a.tui.ShowParameterForm(parameters, ui.OkButtonExecute); ok {
		executeScript(snippet.Format(parameterValues, formatOptions(a.config.Script)), a.config.Script.Shell)
	}
}

func executeScript(script, configuredShell string) {
	shell := stringutil.FirstNotEmpty(configuredShell, os.Getenv("SHELL"), fallbackShell)

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

func formatOptions(cfg config.ScriptConfig) model.SnippetFormatOptions {
	var paramMode model.SnippetParamMode
	if cfg.ParameterMode == config.ParameterModeReplace {
		paramMode = model.SnippetParamModeReplace
	} else {
		paramMode = model.SnippetParamModeSet
	}
	return model.SnippetFormatOptions{
		RemoveComments: cfg.RemoveComments,
		ParamMode:      paramMode,
	}
}
