package chat

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/afero"

	appModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/form"
	"github.com/lemoony/snipkit/internal/ui/style"
)

type focusArea int

const (
	focusFields focusArea = iota
	focusButtons

	// UI constants.
	keyLeft             = "left"
	keyRight            = "right"
	modalContentPadding = 4
)

type parameterModal struct {
	fields       []*form.FieldModel
	elementFocus int
	parameters   []appModel.Parameter
	styler       style.Style
	fs           afero.Fs

	// Focus state
	focusArea   focusArea
	buttonFocus int // 0=Execute, 1=Cancel

	// Tracking state
	submitted bool
	canceled  bool
}

// NewParameterModal creates a new modal for collecting parameter values.
func NewParameterModal(parameters []appModel.Parameter, styler style.Style, fs afero.Fs) *parameterModal {
	fields := make([]*form.FieldModel, len(parameters))
	maxLabelWidth := 0

	// Create fields for each parameter
	for i, param := range parameters {
		fields[i] = form.NewField(
			styler,
			param.Name,
			param.Description,
			param.Type,
			param.Values,
			fs,
		)

		// Pre-fill default value if present
		if param.DefaultValue != "" {
			fields[i].SetValue(param.DefaultValue)
		}

		labelWidth := lipgloss.Width(param.Name)
		if labelWidth > maxLabelWidth {
			maxLabelWidth = labelWidth
		}
	}

	// Set uniform label width for all fields
	for _, field := range fields {
		field.SetLabelWidth(maxLabelWidth)
	}

	return &parameterModal{
		fields:       fields,
		elementFocus: 0,
		parameters:   parameters,
		styler:       styler,
		fs:           fs,
	}
}

// Init initializes the modal.
func (m *parameterModal) Init() tea.Cmd {
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
func (m *parameterModal) Update(msg tea.Msg) (*parameterModal, tea.Cmd) {
	var cmd tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		action := handleModalKeyPress(keyMsg, "e", m.focusArea, m.buttonFocus, m.elementFocus, len(m.fields))

		switch action {
		case modalKeyCancel:
			m.canceled = true
			return m, nil
		case modalKeySubmit:
			m.submitted = true
			return m, nil
		case modalKeyNavigateForward:
			cmd = m.navigateForward()
		case modalKeyNavigateBackward:
			cmd = m.navigateBackward()
		case modalKeyNavigateFields:
			cmd = m.navigateFields(false)
		case modalKeyNavigateButtonsForward:
			m.navigateButtonsForward()
			return m, nil
		case modalKeyNavigateButtonsBackward:
			m.navigateButtonsBackward()
			return m, nil
		}
	}

	// Always delegate to field when in fields area (like form.go does)
	// This allows fields to handle arrow keys for dropdown navigation
	if m.focusArea == focusFields && m.elementFocus < len(m.fields) {
		var fieldCmd tea.Cmd
		m.fields[m.elementFocus], fieldCmd = m.fields[m.elementFocus].Update(msg)
		return m, tea.Batch(cmd, fieldCmd)
	}

	return m, cmd
}

// navigateFields moves focus between fields.
func (m *parameterModal) navigateFields(backward bool) tea.Cmd {
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
func (m *parameterModal) navigateForward() tea.Cmd {
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
			m.buttonFocus = 0 // Focus Execute button
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
func (m *parameterModal) navigateBackward() tea.Cmd {
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
func (m *parameterModal) navigateButtonsForward() {
	m.buttonFocus++
	if m.buttonFocus >= 2 {
		m.buttonFocus = 0
	}
}

// navigateButtonsBackward moves left through buttons.
func (m *parameterModal) navigateButtonsBackward() {
	m.buttonFocus--
	if m.buttonFocus < 0 {
		m.buttonFocus = 1
	}
}

// View renders the modal as a centered box.
func (m *parameterModal) View(terminalWidth, terminalHeight int) string {
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
		m.styler.TitleStyle().Render("Script Parameters"),
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
func (m *parameterModal) renderControlBar() string {
	return renderModalControlBar(
		m.styler,
		m.focusArea,
		m.buttonFocus,
		"Execute",
		"E",
		"Cancel",
	)
}

// GetValues returns the collected parameter values.
func (m *parameterModal) GetValues() []string {
	values := make([]string, len(m.fields))
	for i, field := range m.fields {
		values[i] = field.Value()
	}
	return values
}

// IsSubmitted returns true if the user submitted the form.
func (m *parameterModal) IsSubmitted() bool {
	return m.submitted
}

// IsCanceled returns true if the user canceled the form.
func (m *parameterModal) IsCanceled() bool {
	return m.canceled
}
