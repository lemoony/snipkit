package form

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/afero"

	appModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/style"
)

// focusArea represents which area of the modal has focus.
type focusArea int

const (
	focusFields focusArea = iota
	focusButtons

	// UI constants.
	keyLeft             = "left"
	keyRight            = "right"
	modalContentPadding = 4
	buttonCount         = 2
)

// ParameterModalConfig configures the behavior and appearance of a parameter modal.
type ParameterModalConfig struct {
	Title         string // Modal title (e.g., "Script Parameters")
	OkButtonText  string // Primary button text (e.g., "Execute", "Apply")
	OkShortcut    string // Keyboard shortcut for OK button (e.g., "e", "a")
	ShowAllFields bool   // true = show all fields at once, false = progressive display
	EmbeddedMode  bool   // true = embedded in larger UI, false = standalone program
}

// ParameterModal is a reusable component for collecting parameter values.
type ParameterModal struct {
	config ParameterModalConfig
	styler style.Style
	fs     afero.Fs

	// Field management
	fields       []*FieldModel
	parameters   []appModel.Parameter
	elementFocus int
	showFields   int // For progressive display mode

	// Focus state
	focusArea   focusArea
	buttonFocus int // 0=OK, 1=Cancel

	// Result state
	submitted bool
	canceled  bool
}

// createFields creates field models from parameters and values.
func createFields(
	parameters []appModel.Parameter,
	values []appModel.ParameterValue,
	styler style.Style,
	fs afero.Fs,
) ([]*FieldModel, int) {
	fields := make([]*FieldModel, len(parameters))
	maxLabelWidth := 0

	for i, param := range parameters {
		name := param.Key
		if param.Name != "" {
			name = param.Name
		}

		fields[i] = NewField(styler, name, param.Description, param.Type, param.Values, fs)

		// Pre-fill default value if present
		if param.DefaultValue != "" {
			fields[i].SetValue(param.DefaultValue)
		}

		// Apply provided values
		for _, paramValue := range values {
			if paramValue.Key == param.Key {
				fields[i].SetValue(paramValue.Value)
			}
		}

		labelWidth := lipgloss.Width(name)
		if labelWidth > maxLabelWidth {
			maxLabelWidth = labelWidth
		}
	}

	return fields, maxLabelWidth
}

// NewParameterModal creates a new parameter modal with the given configuration.
func NewParameterModal(
	parameters []appModel.Parameter,
	values []appModel.ParameterValue,
	config ParameterModalConfig,
	styler style.Style,
	fs afero.Fs,
) *ParameterModal {
	fields, maxLabelWidth := createFields(parameters, values, styler, fs)

	// Set uniform label width for all fields
	for _, field := range fields {
		field.SetLabelWidth(maxLabelWidth)
	}

	showFields := len(fields) + buttonCount
	if !config.ShowAllFields {
		showFields = 0
	}

	return &ParameterModal{
		config:       config,
		fields:       fields,
		parameters:   parameters,
		elementFocus: 0,
		showFields:   showFields,
		styler:       styler,
		fs:           fs,
		focusArea:    focusFields,
		buttonFocus:  0,
	}
}

// Init initializes the modal.
func (m *ParameterModal) Init() tea.Cmd {
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

// handleKeyAction processes modal key actions and updates state.
func (m *ParameterModal) handleKeyAction(action modalKeyAction) (*ParameterModal, tea.Cmd) {
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
	case modalKeyNavigateUp:
		return m, m.navigateUpFromButtons()
	case modalKeyNavigateToFirstField:
		return m, m.navigateToFirstField()
	}
	return m, nil
}

// Update handles messages for the modal.
func (m *ParameterModal) Update(msg tea.Msg) (*ParameterModal, tea.Cmd) {
	var cmd tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Special handling for Enter key in fields - check if field has option to apply
		if keyMsg.String() == "enter" && m.focusArea == focusFields && m.elementFocus < len(m.fields) {
			if m.fields[m.elementFocus].HasOptionToApply() {
				// Let field handle the enter to apply the option
				var fieldCmd tea.Cmd
				m.fields[m.elementFocus], fieldCmd = m.fields[m.elementFocus].Update(msg)
				return m, fieldCmd
			}
		}

		action := handleModalKeyPress(
			keyMsg,
			m.config.OkShortcut,
			m.focusArea,
			m.buttonFocus,
			m.elementFocus,
			len(m.fields),
		)

		if action != modalKeyNone {
			updatedModal, actionCmd := m.handleKeyAction(action)
			m = updatedModal
			cmd = actionCmd
		}
	}

	// Always delegate to field when in fields area
	// This allows fields to handle arrow keys for dropdown navigation
	if m.focusArea == focusFields && m.elementFocus < len(m.fields) {
		var fieldCmd tea.Cmd
		m.fields[m.elementFocus], fieldCmd = m.fields[m.elementFocus].Update(msg)
		return m, tea.Batch(cmd, fieldCmd)
	}

	return m, cmd
}

// navigateFields moves focus between fields.
func (m *ParameterModal) navigateFields(backward bool) tea.Cmd {
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

	// Update showFields for progressive mode
	if !m.config.ShowAllFields && m.elementFocus > m.showFields {
		m.showFields = m.elementFocus
	}

	// Focus new field
	return m.fields[m.elementFocus].Focus()
}

// navigateForward moves focus forward through fields and buttons.
func (m *ParameterModal) navigateForward() tea.Cmd {
	if m.focusArea == focusFields {
		// In fields
		if m.elementFocus < len(m.fields)-1 {
			// Move to next field
			m.fields[m.elementFocus].Blur()
			m.elementFocus++

			// Update showFields for progressive mode
			if !m.config.ShowAllFields && m.elementFocus > m.showFields {
				m.showFields = m.elementFocus
			}

			return m.fields[m.elementFocus].Focus()
		}

		// Last field - move to buttons
		m.fields[m.elementFocus].Blur()
		m.focusArea = focusButtons
		m.buttonFocus = 0 // Focus OK button

		// Show buttons in progressive mode
		if !m.config.ShowAllFields {
			m.showFields = len(m.fields) + buttonCount
		}

		return nil
	}

	// In buttons - cycle through buttons
	m.buttonFocus++
	if m.buttonFocus >= buttonCount {
		m.buttonFocus = 0
	}
	return nil
}

// navigateBackward moves focus backward through buttons and fields.
func (m *ParameterModal) navigateBackward() tea.Cmd {
	if m.focusArea == focusButtons {
		// In buttons
		if m.buttonFocus > 0 {
			// Move to previous button
			m.buttonFocus--
			return nil
		}

		// First button - move to last field
		if len(m.fields) > 0 {
			m.focusArea = focusFields
			m.elementFocus = len(m.fields) - 1
			return m.fields[m.elementFocus].Focus()
		}
		return nil
	}

	// In fields
	if m.elementFocus > 0 {
		// Move to previous field
		m.fields[m.elementFocus].Blur()
		m.elementFocus--
		return m.fields[m.elementFocus].Focus()
	}

	// First field - wrap to last button
	m.fields[m.elementFocus].Blur()
	m.focusArea = focusButtons
	m.buttonFocus = 1 // Focus Cancel button
	return nil
}

// navigateButtonsForward moves right through buttons.
func (m *ParameterModal) navigateButtonsForward() {
	m.buttonFocus++
	if m.buttonFocus >= buttonCount {
		m.buttonFocus = 0
	}
}

// navigateButtonsBackward moves left through buttons.
func (m *ParameterModal) navigateButtonsBackward() {
	m.buttonFocus--
	if m.buttonFocus < 0 {
		m.buttonFocus = 1
	}
}

// navigateUpFromButtons moves from buttons to last field.
func (m *ParameterModal) navigateUpFromButtons() tea.Cmd {
	if m.focusArea != focusButtons || len(m.fields) == 0 {
		return nil
	}

	m.focusArea = focusFields
	m.elementFocus = len(m.fields) - 1
	return m.fields[m.elementFocus].Focus()
}

// navigateToFirstField moves from last button to first field (Tab cycling).
func (m *ParameterModal) navigateToFirstField() tea.Cmd {
	if len(m.fields) == 0 {
		return nil
	}

	m.focusArea = focusFields
	m.elementFocus = 0
	return m.fields[0].Focus()
}

// View renders the modal.
func (m *ParameterModal) View(terminalWidth, terminalHeight int) string {
	// Build form content
	var formFields []string
	maxFieldIndex := len(m.fields) - 1
	if !m.config.ShowAllFields {
		maxFieldIndex = m.showFields
		if maxFieldIndex > len(m.fields)-1 {
			maxFieldIndex = len(m.fields) - 1
		}
	}

	for i := 0; i <= maxFieldIndex && i < len(m.fields); i++ {
		formFields = append(formFields, m.fields[i].View())
	}
	formContent := lipgloss.JoinVertical(lipgloss.Left, formFields...)

	// Build content sections
	var contentSections []string
	contentSections = append(contentSections, m.styler.TitleStyle().Render(m.config.Title))

	if len(formFields) > 0 {
		contentSections = append(contentSections, "", formContent)
	}

	// Render control bar if we're showing buttons
	if m.config.ShowAllFields || m.showFields >= len(m.fields) {
		controlBar := m.renderControlBar()
		contentSections = append(contentSections, "", controlBar)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, contentSections...)

	// For embedded mode, return as modal box
	if m.config.EmbeddedMode {
		contentWidth := lipgloss.Width(content) + modalContentPadding

		modalStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.styler.BorderColor().Value()).
			Padding(1, 2).
			Width(contentWidth)

		return modalStyle.Render(content)
	}

	// For standalone mode, return content directly (caller will handle styling)
	return content
}

// renderControlBar renders the interactive button menu.
func (m *ParameterModal) renderControlBar() string {
	return renderModalControlBar(
		m.styler,
		m.focusArea,
		m.buttonFocus,
		m.config.OkButtonText,
		m.config.OkShortcut,
		"Cancel",
	)
}

// GetValues returns the collected parameter values.
func (m *ParameterModal) GetValues() []string {
	values := make([]string, len(m.fields))
	for i, field := range m.fields {
		values[i] = field.Value()
	}
	return values
}

// IsSubmitted returns true if the user submitted the form.
func (m *ParameterModal) IsSubmitted() bool {
	return m.submitted
}

// IsCanceled returns true if the user canceled the form.
func (m *ParameterModal) IsCanceled() bool {
	return m.canceled
}

// GetFocusArea returns the current focus area (for testing).
func (m *ParameterModal) GetFocusArea() int {
	return int(m.focusArea)
}

// GetElementFocus returns the current element focus index (for testing).
func (m *ParameterModal) GetElementFocus() int {
	return m.elementFocus
}

// GetButtonFocus returns the current button focus index (for testing).
func (m *ParameterModal) GetButtonFocus() int {
	return m.buttonFocus
}

// SetFocusArea sets the focus area (for testing).
func (m *ParameterModal) SetFocusArea(area int) {
	m.focusArea = focusArea(area)
}

// SetElementFocus sets the element focus index (for testing).
func (m *ParameterModal) SetElementFocus(focus int) {
	m.elementFocus = focus
}

// SetButtonFocus sets the button focus index (for testing).
func (m *ParameterModal) SetButtonFocus(focus int) {
	m.buttonFocus = focus
}

// renderModalControlBar renders a two-button control bar with the given labels and shortcuts.
func renderModalControlBar(
	styler style.Style,
	focusArea focusArea,
	buttonFocus int,
	primaryLabel string,
	primaryShortcut string,
	secondaryLabel string,
) string {
	buttons := []string{}

	// Primary button (Execute/Apply)
	primaryStyle := lipgloss.NewStyle().Padding(0, 1)
	if focusArea == focusButtons && buttonFocus == 0 {
		primaryStyle = primaryStyle.
			Bold(true).
			Foreground(styler.HighlightColor().Value()).
			Background(styler.ActiveColor().Value())
	} else {
		primaryStyle = primaryStyle.Foreground(styler.TextColor().Value())
	}
	buttons = append(buttons, primaryStyle.Render(primaryLabel))

	// Secondary button (Cancel)
	secondaryStyle := lipgloss.NewStyle().Padding(0, 1)
	if focusArea == focusButtons && buttonFocus == 1 {
		secondaryStyle = secondaryStyle.
			Bold(true).
			Foreground(styler.HighlightColor().Value()).
			Background(styler.ActiveColor().Value())
	} else {
		secondaryStyle = secondaryStyle.Foreground(styler.TextColor().Value())
	}
	buttons = append(buttons, secondaryStyle.Render(secondaryLabel))

	// Join buttons horizontally
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)

	// Add help text below with Ctrl shortcuts
	helpText := lipgloss.NewStyle().
		Foreground(styler.PlaceholderColor().Value()).
		Render("Ctrl+" + strings.ToUpper(primaryShortcut) + ": " + strings.ToLower(primaryLabel) + " • Esc: cancel • Tab/↑/↓: navigate • Enter: select")

	return lipgloss.JoinVertical(lipgloss.Left, buttonRow, helpText)
}

// modalKeyAction represents the result of handling a key in a modal.
type modalKeyAction int

const (
	modalKeyNone modalKeyAction = iota
	modalKeySubmit
	modalKeyCancel
	modalKeyNavigateForward
	modalKeyNavigateBackward
	modalKeyNavigateFields
	modalKeyNavigateButtonsForward
	modalKeyNavigateButtonsBackward
	modalKeyDelegateToField
	modalKeyNavigateUp
	modalKeyNavigateToFirstField
)

// handleEnterKey processes the enter key based on focus area and position.
func handleEnterKey(focusArea focusArea, buttonFocus int, elementFocus int, fieldCount int) modalKeyAction {
	if focusArea == focusButtons {
		if buttonFocus == 0 {
			return modalKeySubmit
		}
		return modalKeyCancel
	}
	// In field - check if it's the last field
	if elementFocus == fieldCount-1 {
		return modalKeyNavigateForward
	}
	return modalKeyNavigateFields
}

// handleTabKey processes the tab key based on focus area and button position.
func handleTabKey(focusArea focusArea, buttonFocus int) modalKeyAction {
	if focusArea == focusButtons && buttonFocus == 1 {
		// Last button - cycle to first field
		return modalKeyNavigateToFirstField
	}
	return modalKeyNavigateForward
}

// handleArrowKeys processes arrow keys based on focus area and key pressed.
func handleArrowKeys(keyStr string, focusArea focusArea) modalKeyAction {
	if keyStr == "up" && focusArea == focusButtons {
		return modalKeyNavigateUp
	}
	if focusArea == focusButtons && (keyStr == keyLeft || keyStr == keyRight) {
		if keyStr == keyLeft {
			return modalKeyNavigateButtonsBackward
		}
		return modalKeyNavigateButtonsForward
	}
	return modalKeyNone
}

// handleModalKeyPress processes common modal key presses and returns the appropriate action.
func handleModalKeyPress(
	keyMsg tea.KeyMsg,
	primaryShortcut string,
	focusArea focusArea,
	buttonFocus int,
	elementFocus int,
	fieldCount int,
) modalKeyAction {
	keyStr := keyMsg.String()

	switch keyStr {
	case "esc":
		return modalKeyCancel
	case "ctrl+" + primaryShortcut:
		return modalKeySubmit
	case "enter":
		return handleEnterKey(focusArea, buttonFocus, elementFocus, fieldCount)
	case "tab":
		return handleTabKey(focusArea, buttonFocus)
	case "shift+tab":
		return modalKeyNavigateBackward
	case "up", keyLeft, keyRight:
		if action := handleArrowKeys(keyStr, focusArea); action != modalKeyNone {
			return action
		}
	}

	// Delegate to field if in field area
	if focusArea == focusFields && elementFocus < fieldCount {
		return modalKeyDelegateToField
	}

	return modalKeyNone
}
