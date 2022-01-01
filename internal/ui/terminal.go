package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/gdamore/tcell/v2"
	"github.com/kballard/go-shellquote"

	"github.com/lemoony/snippet-kit/internal/model"
)

const (
	envEditor            = "EDITOR"
	envVisual            = "VISUAL"
	defaultEditor        = "vim"
	defaultEditorWindows = "notepad"
	windows              = "windows"
)

// TerminalOption configures a Terminal.
type TerminalOption interface {
	apply(p *cliTerminal)
}

// terminalOptionFunc wraps a func so that it satisfies the Option interface.
type terminalOptionFunc func(terminal *cliTerminal)

func (f terminalOptionFunc) apply(terminal *cliTerminal) {
	f(terminal)
}

// WithStdio sets the stdio for the terminal.
func WithStdio(stdio terminal.Stdio) TerminalOption {
	return terminalOptionFunc(func(t *cliTerminal) {
		t.stdio = stdio
	})
}

// WithScreen sets the screen for tview.
func WithScreen(screen tcell.Screen) TerminalOption {
	return terminalOptionFunc(func(t *cliTerminal) {
		t.screen = screen
	})
}

type Terminal interface {
	PrintMessage(message string)
	PrintError(message string)
	Confirm(message string) (bool, error)
	OpenEditor(path string, preferredEditor string) error
	ShowLookup(snippets []model.Snippet) (int, error)
	ShowParameterForm(parameters []model.Parameter) ([]string, error)
}

type cliTerminal struct {
	stdio  terminal.Stdio
	screen tcell.Screen
}

func NewTerminal(options ...TerminalOption) Terminal {
	term := cliTerminal{
		stdio: terminal.Stdio{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
	}
	for _, option := range options {
		option.apply(&term)
	}
	return term
}

func (c cliTerminal) PrintMessage(msg string) {
	fmt.Fprintln(c.stdio.Out, msg)
}

func (c cliTerminal) PrintError(msg string) {
	fmt.Fprintln(c.stdio.Out, msg)
}

func (c cliTerminal) Confirm(message string) (bool, error) {
	confirmed := false
	prompt := &survey.Confirm{Message: message}
	if err := survey.AskOne(prompt, &confirmed, survey.WithStdio(c.stdio.In, c.stdio.Out, c.stdio.Err)); err != nil {
		return false, err
	}
	return confirmed, nil
}

func (c cliTerminal) OpenEditor(path string, preferredEditor string) error {
	editor := getEditor(preferredEditor)

	args, err := shellquote.Split(editor)
	if err != nil {
		return err
	}
	args = append(args, path)

	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec // subprocess launched with c potential tainted input
	cmd.Stdin = c.stdio.In
	cmd.Stdout = c.stdio.Out
	cmd.Stderr = c.stdio.Err

	err = cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func getEditor(preferred string) string {
	result := defaultEditor
	if runtime.GOOS == windows {
		result = defaultEditorWindows
	}

	preferred = strings.TrimSpace(preferred)
	if preferred != "" {
		result = preferred
	} else if v := os.Getenv(envVisual); v != "" {
		result = v
	} else if e := os.Getenv(envEditor); e != "" {
		result = e
	}

	return result
}
