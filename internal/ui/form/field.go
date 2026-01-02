package form

import (
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/afero"

	appModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

const (
	fieldTotalWidth = 75
	fieldMaxOptions = 4

	two                   = 2
	fieldLabelPaddingLeft = 4
)

// FieldModel represents a single parameter input field.
type FieldModel struct {
	styler style.Style
	fs     afero.Fs

	Label           string
	Description     string
	ParameterType   appModel.ParameterType
	field           textinput.Model
	options         []string
	filteredOptions []string
	filterMatches   []int

	selectedPathSuggestion string

	keyMap FieldKeyMap

	labelWidth int

	selectedOption int
	optionOffset   int
}

func NewField(
	styler style.Style,
	label,
	description string,
	paramType appModel.ParameterType,
	options []string,
	fs afero.Fs,
) *FieldModel {
	m := FieldModel{
		styler:         styler,
		fs:             fs,
		Label:          label,
		ParameterType:  paramType,
		Description:    description,
		keyMap:         defaultFieldKeyMap(),
		field:          textinput.New(),
		options:        options,
		selectedOption: -1,
		optionOffset:   0,
	}

	m.field.Prompt = ""
	m.field.Placeholder = stringutil.StringOrDefault(description, "Type here...")
	m.field.Cursor.SetMode(cursor.CursorBlink)
	if m.ParameterType == appModel.ParameterTypePassword {
		m.field.EchoMode = textinput.EchoPassword
	}

	return &m
}

func (m *FieldModel) Value() string {
	return m.field.Value()
}

func (m *FieldModel) SetLabelWidth(width int) {
	m.field.Width = fieldTotalWidth - width - two - 8 //nolint:mnd // magic number 8
	m.labelWidth = width
}

func (m *FieldModel) SetValue(text string) {
	m.field.SetValue(text)
}

func (m *FieldModel) Focus() tea.Cmd {
	// Always filter options on focus to initialize suggestions
	m.filterOptions()
	return m.field.Focus()
}

func (m *FieldModel) Blur() {
	m.field.Blur()
}

func (m *FieldModel) HasOptionToApply() bool {
	return m.selectedPathSuggestion != ""
}

func (m *FieldModel) Update(msg tea.Msg) (*FieldModel, tea.Cmd) {
	var cmd tea.Cmd

	handled := false
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.CursorDown):
			handled = true
			m.selectNextOption()
		case key.Matches(msg, m.keyMap.CursorUp):
			handled = true
			m.selectPreviousOption()
		case key.Matches(msg, m.keyMap.Apply):
			if m.HasOptionToApply() {
				m.applyFilePathOption()
				handled = true
			}
		case key.Matches(msg, m.keyMap.ApplyCompletion):
			// Only handle if we have a suggestion to apply
			if m.HasOptionToApply() {
				m.applyFilePathOption()
				handled = true
			}
			// Otherwise falls through to textinput for normal cursor movement
		}
	}

	// Always update textinput for cursor, but only check value changes if not handled
	prevValue := m.field.Value()
	m.field, cmd = m.field.Update(msg)
	if !handled && prevValue != m.field.Value() {
		m.filterOptions()
	}

	return m, cmd
}

func (m *FieldModel) selectNextOption() {
	if len(m.options) == 0 {
		return
	}

	m.selectedOption = min(m.selectedOption+1, len(m.options)-1)
	if m.selectedOption >= m.optionOffset+fieldMaxOptions {
		m.optionOffset += 1
	}

	if m.ParameterType == appModel.ParameterTypeValue {
		m.field.SetValue(m.options[m.selectedOption])
		m.field.CursorEnd()
		m.filterOptions()
	} else {
		m.selectedPathSuggestion = m.options[m.selectedOption]
	}
}

func (m *FieldModel) selectPreviousOption() {
	if len(m.options) == 0 {
		return
	}

	m.selectedOption = max(m.selectedOption-1, 0)
	if m.selectedOption-m.optionOffset < 0 {
		m.optionOffset -= 1
	}

	if m.ParameterType == appModel.ParameterTypeValue {
		m.field.SetValue(m.options[m.selectedOption])
		m.field.CursorEnd()
		m.filterOptions()
	} else {
		m.selectedPathSuggestion = m.options[m.selectedOption]
	}
}

func (m *FieldModel) filterOptions() {
	if m.ParameterType == appModel.ParameterTypeValue {
		m.filterOptionsForValue()
	} else if m.ParameterType == appModel.ParameterTypePath {
		m.filterOptionsForFilePath()
	}
}

func (m *FieldModel) filterOptionsForValue() {
	filterValue := strings.ToLower(m.field.Value())

	m.filteredOptions = m.options

	var filtered []int
	matchIndices := make([]int, len(m.filteredOptions))

	for i, o := range m.options {
		if foundIndex := strings.Index(strings.ToLower(o), filterValue); foundIndex >= 0 && len(filterValue) > 0 {
			filtered = append(filtered, i)
			matchIndices[i] = foundIndex
		} else {
			matchIndices[i] = -1
		}
	}

	if len(filtered) > 0 {
		m.selectedOption = filtered[0]
	}

	m.filterMatches = matchIndices
	m.selectedPathSuggestion = ""
}

func (m *FieldModel) filterOptionsForFilePath() {
	m.filteredOptions = suggestionsForPath(m.fs, m.field.Value())
	m.options = m.filteredOptions
	m.filterMatches = make([]int, len(m.filteredOptions))
	for i := range m.filteredOptions {
		m.filterMatches[i] = 0
	}

	// Auto-select first option and set it as the proposed suggestion
	if len(m.options) > 0 {
		m.selectedOption = 0
		m.optionOffset = 0
		m.selectedPathSuggestion = m.options[0]
	} else {
		m.selectedOption = -1
		m.selectedPathSuggestion = ""
	}
}

func (m *FieldModel) applyFilePathOption() {
	m.field.SetValue(m.selectedPathSuggestion)
	m.selectedPathSuggestion = ""
	m.field.CursorEnd()
	m.filterOptions()
}

//nolint:funlen // refactor at a later point.
func (m *FieldModel) View() string {
	color := m.styler.SubduedColor()
	borderStyle := lipgloss.HiddenBorder()
	if m.field.Focused() {
		color = m.styler.ActiveColor()
		borderStyle = lipgloss.NormalBorder()
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(color.Value()).
		Bold(true).
		Border(borderStyle, false, false, false, true).
		BorderForeground(color.Value()).
		Padding(0, fieldLabelPaddingLeft, 0, 1)

	label := labelStyle.Render(lipgloss.PlaceHorizontal(m.labelWidth, lipgloss.Left, m.Label, lipgloss.WithWhitespaceChars(" ")))

	var fieldView string
	if m.selectedPathSuggestion != "" && m.ParameterType == appModel.ParameterTypePath {
		// Show typed text in normal style
		typedText := m.field.TextStyle.Render(m.field.Value())

		// Show completion in subdued/grayed style (like zsh)
		completion := strings.TrimPrefix(m.selectedPathSuggestion, m.field.Value())
		completionStyle := lipgloss.NewStyle().
			Foreground(m.styler.SubduedColor().Value())

		fieldView = typedText + completionStyle.Render(completion)
	} else {
		fieldView = m.field.View()
	}

	f := lipgloss.JoinHorizontal(lipgloss.Left, label, fieldView)

	var options string
	if m.field.Focused() {
		for i, o := range m.filteredOptions {
			if i < m.optionOffset {
				continue
			} else if i-m.optionOffset >= fieldMaxOptions {
				options = lipgloss.JoinVertical(lipgloss.Left, options, lipgloss.NewStyle().MarginLeft(lipgloss.Width(label)).Render("..."))
				break
			}

			if len(m.filterMatches) > i {
				if match := m.filterMatches[i]; match >= 0 {
					first := o[0:match]
					middle := o[match : match+len(m.field.Value())]
					end := o[match+len(m.field.Value()):]

					o = lipgloss.JoinHorizontal(lipgloss.Left, first, lipgloss.NewStyle().Foreground(m.styler.ActiveColor().Value()).Render(middle), end)
				}
			}

			x := lipgloss.PlaceHorizontal(lipgloss.Width(label)-two, lipgloss.Left, "")

			if i == m.selectedOption {
				x += "> " + o
			} else {
				x += "  " + o
			}

			if options == "" {
				options = x
			} else {
				options = lipgloss.JoinVertical(lipgloss.Left, options, x)
			}
		}
	}

	if options != "" {
		f = lipgloss.JoinVertical(lipgloss.Left, f, options)
	}

	return f
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
