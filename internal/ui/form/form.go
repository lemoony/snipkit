package form

import (
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"

	internalModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/style"
)

// standaloneWrapper wraps ParameterModal for standalone execution.
type standaloneWrapper struct {
	modal        *ParameterModal
	colorProfile termenv.Profile
	input        *io.Reader
	output       *io.Writer
	help         help.Model
	styler       style.Style
	width        int
	height       int
}

// Show displays a parameter collection form and returns the collected values.
func Show(parameters []internalModel.Parameter, values []internalModel.ParameterValue, okButton string, options ...Option) ([]string, bool) {
	// Apply options to get configuration
	cfg := &config{
		colorProfile: termenv.ColorProfile(),
		help:         help.New(),
	}

	for _, o := range options {
		o.apply(cfg)
	}

	// Create parameter modal with standalone configuration
	modal := NewParameterModal(
		parameters,
		values,
		ParameterModalConfig{
			Title:         "This snippet requires parameters",
			OkButtonText:  okButton,
			OkShortcut:    "a", // Ctrl+A for Apply
			ShowAllFields: false,
			EmbeddedMode:  false,
		},
		cfg.styler,
		cfg.fs,
	)

	// Create standalone wrapper
	wrapper := &standaloneWrapper{
		modal:        modal,
		colorProfile: cfg.colorProfile,
		input:        cfg.input,
		output:       cfg.output,
		help:         cfg.help,
		styler:       cfg.styler,
	}

	// Configure tea program options
	var teaOptions []tea.ProgramOption
	if wrapper.input != nil {
		teaOptions = append(teaOptions, tea.WithInput(*wrapper.input))
	}
	if wrapper.output != nil {
		teaOptions = append(teaOptions, tea.WithOutput(*wrapper.output))
	}
	teaOptions = append(teaOptions, tea.WithAltScreen())

	// Run the program
	p := tea.NewProgram(wrapper, teaOptions...)
	if _, err := p.Run(); err != nil {
		panic(err)
	}

	// Return results - empty values if canceled
	if modal.IsSubmitted() {
		return modal.GetValues(), true
	}
	return []string{}, false
}

// Init initializes the standalone wrapper.
func (w *standaloneWrapper) Init() tea.Cmd {
	return w.modal.Init()
}

// Update handles messages for the standalone wrapper.
func (w *standaloneWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit keys
		if msg.String() == "ctrl+c" {
			return w, tea.Quit
		}

	case tea.WindowSizeMsg:
		w.width = msg.Width
		w.height = msg.Height
		w.help.Width = msg.Width
		w.styler.SetSize(msg.Width, msg.Height)
	}

	// Update modal
	var cmd tea.Cmd
	w.modal, cmd = w.modal.Update(msg)

	// Check if modal is done
	if w.modal.IsSubmitted() || w.modal.IsCanceled() {
		return w, tea.Quit
	}

	return w, cmd
}

// View renders the standalone wrapper.
func (w *standaloneWrapper) View() string {
	modalContent := w.modal.View(w.width, w.height)
	helpView := w.help.View(w)

	// Use styler's MainView for consistent full-screen layout
	res := w.styler.MainView(modalContent, helpView, false)
	if w.styler.NeedsResize() {
		res = w.styler.MainView(modalContent, helpView, true)
	}
	return res
}

// FullHelp returns bindings to show the full help view. It's part of the
// help.KeyMap interface.
func (w *standaloneWrapper) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

// ShortHelp returns bindings to show in the abbreviated help view. It's part
// of the help.KeyMap interface.
func (w *standaloneWrapper) ShortHelp() []key.Binding {
	km := defaultKeyMap()
	return []key.Binding{
		km.Next,
		km.Quit,
		km.Apply,
	}
}
