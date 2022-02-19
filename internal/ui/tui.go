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
	"github.com/phuslu/log"
	"github.com/rivo/tview"
	"github.com/spf13/afero"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/confirm"
	"github.com/lemoony/snipkit/internal/ui/form"
	"github.com/lemoony/snipkit/internal/ui/picker"
	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/ui/sync"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/termutil"
)

type OkButton string

const (
	envEditor     = "EDITOR"
	envVisual     = "VISUAL"
	defaultEditor = "vim"

	OkButtonExecute = OkButton("Execute")
	OkButtonPrint   = OkButton("Print")
)

// TUIOption configures a TUI.
type TUIOption interface {
	apply(p *tuiImpl)
}

// tuiOptionFunc wraps a func so that it satisfies the Option interface.
type tuiOptionFunc func(terminal *tuiImpl)

func (f tuiOptionFunc) apply(terminal *tuiImpl) {
	f(terminal)
}

// WithStdio sets the stdio for the terminal.
func WithStdio(stdio termutil.Stdio) TUIOption {
	return tuiOptionFunc(func(t *tuiImpl) {
		t.stdio = stdio
	})
}

// WithScreen sets the screen for tview.
func WithScreen(screen tcell.Screen) TUIOption {
	return tuiOptionFunc(func(t *tuiImpl) {
		t.screen = screen
	})
}

type TUI interface {
	ApplyConfig(cfg Config, system *system.System)
	Print(m uimsg.Printable)
	PrintMessage(message string)
	PrintError(message string)
	Confirmation(confirm uimsg.Confirm, options ...confirm.Option) bool
	OpenEditor(path string, preferredEditor string)
	ShowLookup(snippets []model.Snippet) int
	ShowParameterForm(parameters []model.Parameter, okButton OkButton) ([]string, bool)
	ShowPicker(items []picker.Item, options ...tea.ProgramOption) (int, bool)
	ShowSync() sync.Screen
}

type tuiImpl struct {
	stdio  termutil.Stdio
	screen tcell.Screen
	styler style.Style
}

func NewTUI(options ...TUIOption) TUI {
	term := tuiImpl{
		stdio: termutil.Stdio{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
	}
	for _, option := range options {
		option.apply(&term)
	}

	return &term
}

func (t *tuiImpl) ApplyConfig(cfg Config, system *system.System) {
	themeValues := cfg.GetSelectedTheme(system)
	t.styler = style.NewStyle(&themeValues, !cfg.HideKeyMap)

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorReset
	tview.Styles.BorderColor = t.styler.BorderColor().CellValue()
	tview.Styles.TitleColor = t.styler.BorderTitleColor().CellValue()

	log.Trace().Msgf("Color profile: %d", t.styler.Profile())
}

func (t tuiImpl) Print(p uimsg.Printable) {
	fmt.Fprintln(t.stdio.Out, p.RenderWith(&t.styler))
}

func (t tuiImpl) PrintMessage(msg string) {
	fmt.Fprintln(t.stdio.Out, msg)
}

func (t tuiImpl) PrintError(msg string) {
	fmt.Fprintln(t.stdio.Out, msg)
}

func (t tuiImpl) ShowParameterForm(parameters []model.Parameter, okButton OkButton) ([]string, bool) {
	if len(parameters) == 0 {
		return []string{}, true
	}

	return form.Show(parameters,
		string(okButton),
		form.WithStyler(t.styler),
		form.WithIn(t.stdio.In),
		form.WithOut(t.stdio.Out),
		form.WithFS(afero.NewOsFs()),
	)
}

func (t tuiImpl) Confirmation(confirmation uimsg.Confirm, options ...confirm.Option) bool {
	return confirm.Show(
		confirmation,
		append(
			[]confirm.Option{
				confirm.WithStyler(t.styler),
				confirm.WithIn(t.stdio.In),
				confirm.WithOut(t.stdio.Out),
			},
			options...,
		)...,
	)
}

func (t tuiImpl) OpenEditor(path string, preferredEditor string) {
	editor := getEditor(preferredEditor)

	args, err := shellquote.Split(editor)
	if err != nil {
		panic(errors.Wrap(err, "failed to correctly format editor command"))
	}
	args = append(args, path)

	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec // subprocess launched with t potential tainted input
	cmd.Stdin = t.stdio.In
	cmd.Stdout = t.stdio.Out
	cmd.Stderr = t.stdio.Err

	err = cmd.Start()
	if err != nil {
		panic(errors.Wrapf(errors.WithStack(err), "failed to open editor: %s", strings.Join(args, " ")))
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}

func (t tuiImpl) ShowPicker(items []picker.Item, options ...tea.ProgramOption) (int, bool) {
	return picker.ShowPicker(items, &t.styler, append(
		[]tea.ProgramOption{
			tea.WithInput(t.stdio.In),
			tea.WithOutput(t.stdio.Out),
		},
		options...)...,
	)
}

func (t tuiImpl) ShowSync() sync.Screen {
	return sync.New(
		sync.WithOut(t.stdio.Out),
		sync.WithIn(t.stdio.In),
		sync.WithStyler(t.styler),
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
