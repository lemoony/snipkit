package app

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"emperror.dev/errors"
	"github.com/creack/pty"
	"github.com/phuslu/log"
	"golang.org/x/term"

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
	shell := detectShell(script, configuredShell)

	//nolint:gosec // since it would report G204 complaining about using a variable as input for exec.Command
	cmd := exec.Command(shell, "-c", script)

	// Check if stdin is a terminal to decide execution mode
	if term.IsTerminal(int(os.Stdin.Fd())) {
		return executeWithPTY(cmd)
	}
	return executeWithoutPTY(cmd)
}

// executeWithPTY runs the command with a pseudo-terminal to preserve colors and interactivity.
func executeWithPTY(cmd *exec.Cmd) *capturedOutput {
	// Get current terminal size
	rows, cols := 24, 80 // defaults
	if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		cols, rows = w, h
	}

	// Start the command with a pty
	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})
	if err != nil {
		panic(errors.Wrapf(errors.WithStack(err), "failed to start command with pty"))
	}
	defer func() { _ = ptmx.Close() }()

	// Set stdin to raw mode to pass through all input
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err == nil {
		defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()
	}

	// Buffer to capture output while also displaying it
	var outputBuf bytes.Buffer

	// Copy stdin to pty in a goroutine
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()

	// Copy pty output to both stdout and buffer
	_, _ = io.Copy(io.MultiWriter(os.Stdout, &outputBuf), ptmx)

	// Wait for command to complete
	if waitErr := cmd.Wait(); waitErr != nil {
		log.Info().Err(waitErr)
	}

	return &capturedOutput{
		stdout: outputBuf.String(),
		stderr: "", // PTY combines stdout and stderr
	}
}

// executeWithoutPTY runs the command without a PTY (for non-terminal contexts).
func executeWithoutPTY(cmd *exec.Cmd) *capturedOutput {
	// Create buffers to capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer

	// Create MultiWriters to write to both os.Stdout/os.Stderr and capture buffers
	stdoutWriter := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderrWriter := io.MultiWriter(os.Stderr, &stderrBuf)

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

// detectShell determines which shell to use for script execution.
// Priority: shebang in script > configured shell > $SHELL env var > fallback.
func detectShell(script, configuredShell string) string {
	// Check for shebang in the script
	if strings.HasPrefix(script, "#!") {
		if idx := strings.Index(script, "\n"); idx != -1 {
			shebang := strings.TrimSpace(script[2:idx])
			// Handle "#!/usr/bin/env bash" style shebangs
			if strings.HasPrefix(shebang, "/usr/bin/env ") {
				interpreter := strings.TrimSpace(strings.TrimPrefix(shebang, "/usr/bin/env "))
				// Remove any arguments after the interpreter name
				if spaceIdx := strings.Index(interpreter, " "); spaceIdx != -1 {
					interpreter = interpreter[:spaceIdx]
				}
				return interpreter
			}
			// Handle direct path shebangs like "#!/bin/bash" or "#!/bin/zsh"
			// Remove any arguments after the path
			if spaceIdx := strings.Index(shebang, " "); spaceIdx != -1 {
				shebang = shebang[:spaceIdx]
			}
			return shebang
		}
	}

	return stringutil.FirstNotEmpty(configuredShell, os.Getenv("SHELL"), fallbackShell)
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
