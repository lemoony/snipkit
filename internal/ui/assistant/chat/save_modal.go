package chat

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/afero"

	appModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/form"
	"github.com/lemoony/snipkit/internal/ui/style"
)

type saveModal struct {
	fields       []*form.FieldModel
	elementFocus int
	styler       style.Style
	fs           afero.Fs

	// Focus state
	focusArea   focusArea
	buttonFocus int // 0=Save, 1=Cancel

	// Tracking state
	submitted bool
	canceled  bool
}

// NewSaveModal creates a new modal for collecting save filename and snippet name.
func NewSaveModal(proposedFilename, proposedSnippetName string, styler style.Style, fs afero.Fs) *saveModal {
	fields := make([]*form.FieldModel, 2)

	// Create filename field
	fields[0] = form.NewField(
		styler,
		"Filename",
		"The file where the script will be saved",
		appModel.ParameterTypeValue,
		nil,
		fs,
	)
	fields[0].SetValue(proposedFilename)

	// Create snippet name field
	fields[1] = form.NewField(
		styler,
		"Snippet Name",
		"The display name for this snippet",
		appModel.ParameterTypeValue,
		nil,
		fs,
	)
	fields[1].SetValue(proposedSnippetName)

	// Set uniform label width for both fields
	maxLabelWidth := lipgloss.Width("Snippet Name") // Longest label
	for _, field := range fields {
		field.SetLabelWidth(maxLabelWidth)
	}

	return &saveModal{
		fields:       fields,
		elementFocus: 0,
		styler:       styler,
		fs:           fs,
	}
}

// Init initializes the modal.
func (m *saveModal) Init() tea.Cmd {
	if len(m.fields) > 0 {
		// Start with focus on first field
		m.focusArea = focusFields
		m.elementFocus = 0
		m.buttonFocus = 0
		return m.fields[0].Focus()
	}
	// No fields - focus buttons
	m.focusArea = focusButtons
	m.buttonFocus = 0
	return nil
}

// Update handles messages for the modal.
func (m *saveModal) Update(msg tea.Msg) (*saveModal, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		action := handleModalKeyPress(keyMsg, "s", m.focusArea, m.buttonFocus, m.elementFocus, len(m.fields))

		switch action {
		case modalKeyCancel:
			m.canceled = true
			return m, nil
		case modalKeySubmit:
			m.submitted = true
			return m, nil
		case modalKeyNavigateForward:
			return m, m.navigateForward()
		case modalKeyNavigateBackward:
			return m, m.navigateBackward()
		case modalKeyNavigateFields:
			return m, m.navigateFields(false)
		case modalKeyNavigateButtonsForward:
			m.navigateButtonsForward()
			return m, nil
		case modalKeyNavigateButtonsBackward:
			m.navigateButtonsBackward()
			return m, nil
		case modalKeyDelegateToField:
			var cmd tea.Cmd
			m.fields[m.elementFocus], cmd = m.fields[m.elementFocus].Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// navigateFields moves focus between fields.
func (m *saveModal) navigateFields(backward bool) tea.Cmd {
	if len(m.fields) == 0 {
		return nil
	}

	// Blur current field
	m.fields[m.elementFocus].Blur()

	// Move focus
	if backward {
		m.elementFocus--
		if m.elementFocus < 0 {
			m.elementFocus = len(m.fields) - 1
		}
	} else {
		m.elementFocus++
		if m.elementFocus >= len(m.fields) {
			m.elementFocus = 0
		}
	}

	// Focus new field
	return m.fields[m.elementFocus].Focus()
}

// navigateForward moves focus forward through fields and buttons.
func (m *saveModal) navigateForward() tea.Cmd {
	if m.focusArea == focusFields {
		// In fields
		if m.elementFocus < len(m.fields)-1 {
			// Move to next field
			m.fields[m.elementFocus].Blur()
			m.elementFocus++
			return m.fields[m.elementFocus].Focus()
		} else {
			// Last field - move to buttons
			m.fields[m.elementFocus].Blur()
			m.focusArea = focusButtons
			m.buttonFocus = 0 // Focus Save button
			return nil
		}
	} else {
		// In buttons - cycle through buttons
		m.buttonFocus++
		if m.buttonFocus >= 2 {
			m.buttonFocus = 0
		}
		return nil
	}
}

// navigateBackward moves focus backward through buttons and fields.
func (m *saveModal) navigateBackward() tea.Cmd {
	if m.focusArea == focusButtons {
		// In buttons
		if m.buttonFocus > 0 {
			// Move to previous button
			m.buttonFocus--
			return nil
		} else {
			// First button - move to last field
			m.focusArea = focusFields
			m.elementFocus = len(m.fields) - 1
			return m.fields[m.elementFocus].Focus()
		}
	} else {
		// In fields
		if m.elementFocus > 0 {
			// Move to previous field
			m.fields[m.elementFocus].Blur()
			m.elementFocus--
			return m.fields[m.elementFocus].Focus()
		} else {
			// First field - wrap to last button
			m.fields[m.elementFocus].Blur()
			m.focusArea = focusButtons
			m.buttonFocus = 1 // Focus Cancel button
			return nil
		}
	}
}

// navigateButtonsForward moves right through buttons.
func (m *saveModal) navigateButtonsForward() {
	m.buttonFocus++
	if m.buttonFocus >= 2 {
		m.buttonFocus = 0
	}
}

// navigateButtonsBackward moves left through buttons.
func (m *saveModal) navigateButtonsBackward() {
	m.buttonFocus--
	if m.buttonFocus < 0 {
		m.buttonFocus = 1
	}
}

// View renders the modal as a centered box.
func (m *saveModal) View(terminalWidth, terminalHeight int) string {
	// Build form content
	var formFields []string
	for _, field := range m.fields {
		formFields = append(formFields, field.View())
	}
	formContent := lipgloss.JoinVertical(lipgloss.Left, formFields...)

	// Render control bar
	controlBar := m.renderControlBar()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styler.TitleStyle().Render("Save Snippet"),
		"",
		formContent,
		"",
		controlBar,
	)

	// Calculate modal dimensions
	contentWidth := lipgloss.Width(content) + modalContentPadding

	// Create modal box with border
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styler.BorderColor().Value()).
		Padding(1, 2).
		Width(contentWidth)

	modal := modalStyle.Render(content)

	// Return the modal without positioning - preview will handle overlay positioning
	return modal
}

// renderControlBar renders the interactive button menu.
func (m *saveModal) renderControlBar() string {
	return renderModalControlBar(
		m.styler,
		m.focusArea,
		m.buttonFocus,
		"Save",
		"S",
		"Cancel",
	)
}

// GetFilename returns the filename value.
func (m *saveModal) GetFilename() string {
	if len(m.fields) > 0 {
		return m.fields[0].Value()
	}
	return ""
}

// GetSnippetName returns the snippet name value.
func (m *saveModal) GetSnippetName() string {
	if len(m.fields) > 1 {
		return m.fields[1].Value()
	}
	return ""
}

// IsSubmitted returns true if the user submitted the form.
func (m *saveModal) IsSubmitted() bool {
	return m.submitted
}

// IsCanceled returns true if the user canceled the form.
func (m *saveModal) IsCanceled() bool {
	return m.canceled
}
