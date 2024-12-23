package app

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"emperror.dev/errors"
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

const fallbackShell = "/bin/bash"

type capturedOutput struct {
	stdout string
	stderr string
}

func (a *appImpl) LookupAndExecuteSnippet(confirm, print bool) {
	if ok, snippet := a.LookupSnippet(); ok {
		parameters := snippet.GetParameters()
		if parameterValues, paramOk := a.tui.ShowParameterForm(parameters, nil, ui.OkButtonExecute); paramOk {
			a.executeSnippet(confirm, print, snippet, parameterValues)
		}
	}
}

func (a *appImpl) FindScriptAndExecuteWithParameters(id string, paramValues []model.ParameterValue, confirm, print bool) {
	if snippetFound, snippet := a.getSnippet(id); !snippetFound {
		panic(ErrSnippetIDNotFound)
	} else if paramOk, parameters := matchParameters(paramValues, snippet.GetParameters()); paramOk {
		a.executeSnippet(confirm, print, snippet, parameters)
	} else if parameterValues, formOk := a.tui.ShowParameterForm(snippet.GetParameters(), paramValues, ui.OkButtonExecute); formOk {
		a.executeSnippet(confirm, print, snippet, parameterValues)
	}
}

func (a *appImpl) getSnippet(id string) (bool, model.Snippet) {
	snippets := a.getAllSnippets()
	for i := range snippets {
		if snippets[i].GetID() == id {
			return true, snippets[i]
		}
	}
	return false, nil
}

func matchParameters(paramValues []model.ParameterValue, snippetParameters []model.Parameter) (bool, []string) {
	result := make([]string, len(snippetParameters))
	found := 0
	for i, parameter := range snippetParameters {
		for _, parameterValue := range paramValues {
			if parameterValue.Key == parameter.Key {
				result[i] = parameterValue.Value
				found++
			}
		}
	}
	return found == len(snippetParameters), result
}

func (a *appImpl) executeSnippet(confirm bool, print bool, snippet model.Snippet, parameterValues []string) (bool, *capturedOutput) {
	script := snippet.Format(parameterValues, formatOptions(a.config.Script))

	if (confirm || a.config.Script.ExecConfirm) && !a.tui.Confirmation(uimsg.ExecConfirm(snippet.GetTitle(), script)) {
		return false, nil
	}

	log.Trace().Msg(script)
	if print || a.config.Script.ExecPrint {
		a.tui.Print(uimsg.ExecPrint(snippet.GetTitle(), script))
	}

	return true, executeScript(script, a.config.Script.Shell)
}

func executeScript(script, configuredShell string) *capturedOutput {
	shell := stringutil.FirstNotEmpty(configuredShell, os.Getenv("SHELL"), fallbackShell)

	// Create buffers to capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer

	// Create MultiWriters to write to both os.Stdout/os.Stderr and capture buffers
	stdoutWriter := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderrWriter := io.MultiWriter(os.Stderr, &stderrBuf)

	//nolint:gosec // since it would report G204 complaining about using a variable as input for exec.Command
	cmd := exec.Command(shell, "-c", script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter

	err := cmd.Start()
	if err != nil {
		panic(errors.Wrapf(errors.WithStack(err), "failed to run command"))
	}

	if err = cmd.Wait(); err != nil {
		log.Info().Err(err)
	}

	return &capturedOutput{
		stdout: stdoutBuf.String(),
		stderr: stderrBuf.String(),
	}
}

func formatOptions(cfg config.ScriptConfig) model.SnippetFormatOptions {
	var paramMode model.SnippetParamMode
	if strings.EqualFold(string(config.ParameterModeReplace), string(cfg.ParameterMode)) {
		paramMode = model.SnippetParamModeReplace
	} else {
		paramMode = model.SnippetParamModeSet
	}
	return model.SnippetFormatOptions{
		RemoveComments: cfg.RemoveComments,
		ParamMode:      paramMode,
	}
}
