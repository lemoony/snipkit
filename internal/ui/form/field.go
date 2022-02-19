package form

import (
	"strings"

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

type fieldModel struct {
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
) *fieldModel {
	m := fieldModel{
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
	m.field.SetCursorMode(textinput.CursorStatic)
	if m.ParameterType == appModel.ParameterTypePassword {
		m.field.EchoMode = textinput.EchoPassword
	}

	return &m
}

func (m *fieldModel) Value() string {
	return m.field.Value()
}

func (m *fieldModel) SetLabelWidth(width int) {
	m.field.Width = fieldTotalWidth - width - two - 8 //nolint:gomnd // TODO refactor at some point
	m.labelWidth = width
}

func (m *fieldModel) SetValue(text string) {
	m.field.SetValue(text)
}

func (m *fieldModel) Focus() tea.Cmd {
	if m.ParameterType != appModel.ParameterTypePath {
		m.filterOptions()
	}
	return m.field.Focus()
}

func (m *fieldModel) Blur() {
	m.field.Blur()
}

func (m *fieldModel) HasOptionToApply() bool {
	return m.selectedPathSuggestion != ""
}

func (m *fieldModel) Update(msg tea.Msg) (*fieldModel, tea.Cmd) {
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
		}
	}

	if !handled {
		prevValue := m.field.Value()
		m.field, cmd = m.field.Update(msg)
		if prevValue != m.field.Value() {
			m.filterOptions()
		}
	}

	return m, cmd
}

func (m *fieldModel) selectNextOption() {
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

func (m *fieldModel) selectPreviousOption() {
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

func (m *fieldModel) filterOptions() {
	if m.ParameterType == appModel.ParameterTypeValue {
		m.filterOptionsForValue()
	} else if m.ParameterType == appModel.ParameterTypePath {
		m.filterOptionsForFilePath()
	}
}

func (m *fieldModel) filterOptionsForValue() {
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

func (m *fieldModel) filterOptionsForFilePath() {
	m.filteredOptions = suggestionsForPath(m.fs, m.field.Value())
	m.options = m.filteredOptions
	m.filterMatches = make([]int, len(m.filteredOptions))
	for i := range m.filteredOptions {
		m.filterMatches[i] = 0
	}
}

func (m *fieldModel) applyFilePathOption() {
	m.field.SetValue(m.selectedPathSuggestion)
	m.selectedPathSuggestion = ""
	m.field.CursorEnd()
	m.filterOptions()
}

//nolint:funlen // refactor at a later point.
func (m *fieldModel) View() string {
	color := m.styler.TextColor()
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
	if m.selectedPathSuggestion != "" {
		fieldView = m.field.TextStyle.Render(m.field.Value()) +
			lipgloss.NewStyle().Italic(true).Foreground(m.styler.PlaceholderColor().Value()).
				Render(strings.TrimPrefix(m.selectedPathSuggestion, m.field.Value()))
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
