package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"emperror.dev/errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gdamore/tcell/v2"
	"github.com/kballard/go-shellquote"
	"github.com/rivo/tview"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/confirm"
	"github.com/lemoony/snipkit/internal/ui/picker"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/termutil"
)

const (
	envEditor     = "EDITOR"
	envVisual     = "VISUAL"
	defaultEditor = "vim"
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
func WithStdio(stdio termutil.Stdio) TerminalOption {
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
	ApplyConfig(cfg Config, system *system.System)
	PrintMessage(message string)
	PrintError(message string)
	Confirmation(confirm uimsg.Confirm, options ...confirm.Option) bool
	OpenEditor(path string, preferredEditor string)
	ShowLookup(snippets []model.Snippet) int
	ShowParameterForm(parameters []model.Parameter, okButton OkButton) ([]string, bool)
	ShowPicker(items []picker.Item, options ...tea.ProgramOption) (int, bool)
}

type cliTerminal struct {
	stdio  termutil.Stdio
	screen tcell.Screen
}

func NewTerminal(options ...TerminalOption) Terminal {
	term := cliTerminal{
		stdio: termutil.Stdio{
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

func (c cliTerminal) ApplyConfig(cfg Config, system *system.System) {
	theme := cfg.GetSelectedTheme(system)

	tview.Styles.PrimitiveBackgroundColor = theme.backgroundColor()
	tview.Styles.BorderColor = theme.borderColor()
	tview.Styles.TitleColor = theme.borderTitleColor()
}

func (c cliTerminal) PrintMessage(msg string) {
	fmt.Fprintln(c.stdio.Out, msg)
}

func (c cliTerminal) PrintError(msg string) {
	fmt.Fprintln(c.stdio.Out, msg)
}

func (c cliTerminal) Confirmation(confirmation uimsg.Confirm, options ...confirm.Option) bool {
	return confirm.Confirm(
		confirmation,
		append(
			[]confirm.Option{
				confirm.WithSelectionColor(currentTheme.PromptSelectionTextColor),
				confirm.WithIn(c.stdio.In),
				confirm.WithOut(c.stdio.Out),
			},
			options...,
		)...,
	)
}

func (c cliTerminal) OpenEditor(path string, preferredEditor string) {
	editor := getEditor(preferredEditor)

	args, err := shellquote.Split(editor)
	if err != nil {
		panic(errors.Wrap(err, "failed to correctly format editor command"))
	}
	args = append(args, path)

	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec // subprocess launched with c potential tainted input
	cmd.Stdin = c.stdio.In
	cmd.Stdout = c.stdio.Out
	cmd.Stderr = c.stdio.Err

	err = cmd.Start()
	if err != nil {
		panic(errors.Wrapf(errors.WithStack(err), "failed to open editor: %s", strings.Join(args, " ")))
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}

func (c cliTerminal) ShowPicker(items []picker.Item, options ...tea.ProgramOption) (int, bool) {
	return picker.ShowPicker(items, append(
		[]tea.ProgramOption{
			tea.WithInput(c.stdio.In),
			tea.WithOutput(c.stdio.Out),
		},
		options...)...,
	)
}

func getEditor(preferred string) string {
	result := defaultEditor

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
