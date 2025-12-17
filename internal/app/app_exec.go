package app

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/phuslu/log"
	"golang.org/x/term"

	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/ui/execution"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

const fallbackShell = "/bin/bash"

// ExecutionContext indicates the origin of the execution request.
type ExecutionContext int

const (
	// ContextDefault represents normal execution (e.g., from lookup).
	ContextDefault ExecutionContext = iota
	// ContextAssistant represents execution initiated from the assistant.
	ContextAssistant
)

type capturedOutput struct {
	stdout   string
	stderr   string
	exitCode int
	duration time.Duration
	err      error
}

// Terminal function variables that can be overridden in tests.
var (
	isTerminalFunc  = term.IsTerminal
	getTermSizeFunc = term.GetSize
	makeRawFunc     = term.MakeRaw
	restoreTermFunc = term.Restore
)

func (a *appImpl) LookupAndExecuteSnippet(confirm, print bool) {
	if ok, snippet := a.LookupSnippet(); ok {
		parameters := snippet.GetParameters()
		if parameterValues, paramOk := a.tui.ShowParameterForm(parameters, nil, ui.OkButtonExecute); paramOk {
			a.executeSnippet(ContextDefault, print, snippet, parameterValues)
		}
	}
}

func (a *appImpl) FindScriptAndExecuteWithParameters(id string, paramValues []model.ParameterValue, confirm, print bool) {
	if snippetFound, snippet := a.getSnippet(id); !snippetFound {
		panic(ErrSnippetIDNotFound)
	} else if paramOk, parameters := matchParameters(paramValues, snippet.GetParameters()); paramOk {
		a.executeSnippet(ContextDefault, print, snippet, parameters)
	} else if parameterValues, formOk := a.tui.ShowParameterForm(snippet.GetParameters(), paramValues, ui.OkButtonExecute); formOk {
		a.executeSnippet(ContextDefault, print, snippet, parameterValues)
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

func (a *appImpl) executeSnippet(context ExecutionContext, print bool, snippet model.Snippet, parameterValues []string) *capturedOutput {
	script := snippet.Format(parameterValues, formatOptions(a.config.Script))

	// Skip confirmation for assistant context (parameter modal serves as implicit confirmation)
	if context == ContextDefault && a.config.Script.ExecConfirm && !a.tui.Confirmation(uimsg.ExecConfirm(snippet.GetTitle(), script)) {
		return nil
	}

	log.Trace().Msg(script)
	if print || a.config.Script.ExecPrint {
		a.tui.Print(uimsg.ExecPrint(snippet.GetTitle(), script))
	}

	return executeScript(script, a.config.Script.Shell)
}

func executeScript(script, configuredShell string) *capturedOutput {
	shell := detectShell(script, configuredShell)

	//nolint:gosec // since it would report G204 complaining about using a variable as input for exec.Command
	cmd := exec.Command(shell, "-c", script)

	// Run the script
	if isTerminalFunc(int(os.Stdin.Fd())) {
		// Use Tea-based viewer for terminal execution
		result := execution.RunWithViewer(cmd)
		return &capturedOutput{
			stdout:   result.Stdout,
			exitCode: result.ExitCode,
			duration: result.Duration,
		}
	}

	return executeWithoutPTY(cmd)
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

	// Track start time
	startTime := time.Now()

	err := cmd.Start()
	if err != nil {
		panic(errors.Wrapf(errors.WithStack(err), "failed to run command"))
	}

	err = cmd.Wait()
	duration := time.Since(startTime)

	// Extract exit code
	exitCode := 0
	if err != nil {
		log.Info().Err(err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1 // Could not determine exit code
		}
	}

	return &capturedOutput{
		stdout:   stdoutBuf.String(),
		stderr:   stderrBuf.String(),
		exitCode: exitCode,
		duration: duration,
		err:      err,
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
