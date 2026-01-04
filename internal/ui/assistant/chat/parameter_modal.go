package chat

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"

	appModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/form"
	"github.com/lemoony/snipkit/internal/ui/style"
)

// parameterModal wraps the common form.ParameterModal for use in the chat UI.
type parameterModal struct {
	modal *form.ParameterModal
}

// NewParameterModal creates a new modal for collecting parameter values.
func NewParameterModal(parameters []appModel.Parameter, styler style.Style, fs afero.Fs) *parameterModal {
	modal := form.NewParameterModal(
		parameters,
		nil, // No pre-filled values in chat context
		form.ParameterModalConfig{
			Title:         "Script Parameters",
			OkButtonText:  "Execute",
			OkShortcut:    "e",
			ShowAllFields: true, // Show all fields at once in chat
			EmbeddedMode:  true, // Embedded in chat UI
		},
		styler,
		fs,
	)

	return &parameterModal{
		modal: modal,
	}
}

// Init initializes the modal.
func (m *parameterModal) Init() tea.Cmd {
	return m.modal.Init()
}

// Update handles messages for the modal.
func (m *parameterModal) Update(msg tea.Msg) (*parameterModal, tea.Cmd) {
	var cmd tea.Cmd
	m.modal, cmd = m.modal.Update(msg)
	return m, cmd
}

// View renders the modal as a centered box.
func (m *parameterModal) View(terminalWidth, terminalHeight int) string {
	return m.modal.View(terminalWidth, terminalHeight)
}

// GetValues returns the collected parameter values.
func (m *parameterModal) GetValues() []string {
	return m.modal.GetValues()
}

// IsSubmitted returns true if the user submitted the form.
func (m *parameterModal) IsSubmitted() bool {
	return m.modal.IsSubmitted()
}

// IsCanceled returns true if the user canceled the form.
func (m *parameterModal) IsCanceled() bool {
	return m.modal.IsCanceled()
}
