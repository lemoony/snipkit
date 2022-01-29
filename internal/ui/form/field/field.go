package field

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

const (
	totalWidth = 75
	maxOptions = 4

	two              = 2
	labelPaddingLeft = 4
)

type Model struct {
	styler style.Style

	Label           string
	Description     string
	field           textinput.Model
	options         []string
	filteredOptions []string
	filterMatches   []int

	keyMap KeyMap

	labelWidth int

	selectedOption int
	optionOffset   int
}

func New(styler style.Style, label, description string, options []string) *Model {
	m := Model{
		styler:         styler,
		Label:          label,
		Description:    description,
		keyMap:         defaultKeyMap(),
		field:          textinput.New(),
		options:        options,
		selectedOption: -1,
		optionOffset:   0,
	}

	m.field.Prompt = ""
	m.field.Placeholder = stringutil.StringOrDefault(description, "Type here...")
	m.field.SetCursorMode(textinput.CursorStatic)
	return &m
}

func (m *Model) Value() string {
	return m.field.Value()
}

func (m *Model) SetLabelWidth(width int) {
	m.field.Width = totalWidth - width - two - 8 //nolint:gomnd // TODO refactor at some point
	m.labelWidth = width
}

func (m *Model) SetValue(text string) {
	m.field.SetValue(text)
}

func (m *Model) Focus() tea.Cmd {
	m.filterOptions()
	return m.field.Focus()
}

func (m *Model) Blur() {
	m.field.Blur()
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
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

func (m *Model) selectNextOption() {
	if len(m.options) == 0 {
		return
	}

	m.selectedOption = min(m.selectedOption+1, len(m.options)-1)
	if m.selectedOption >= m.optionOffset+maxOptions {
		m.optionOffset += 1
	}

	m.field.SetValue(m.options[m.selectedOption])
	m.field.CursorEnd()
	m.filterOptions()
}

func (m *Model) selectPreviousOption() {
	if len(m.options) == 0 {
		return
	}

	m.selectedOption = max(m.selectedOption-1, 0)
	if m.selectedOption-m.optionOffset < 0 {
		m.optionOffset -= 1
	}

	m.field.SetValue(m.options[m.selectedOption])
	m.field.CursorEnd()
	m.filterOptions()
}

func (m *Model) filterOptions() {
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
}

func (m *Model) View() string {
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
		Padding(0, labelPaddingLeft, 0, 1)

	label := labelStyle.Render(lipgloss.PlaceHorizontal(m.labelWidth, lipgloss.Left, m.Label, lipgloss.WithWhitespaceChars(" ")))

	f := lipgloss.JoinHorizontal(lipgloss.Left, label, m.field.View())

	var options string
	if m.field.Focused() {
		for i, o := range m.filteredOptions {
			if i < m.optionOffset {
				continue
			} else if i-m.optionOffset >= maxOptions {
				options = lipgloss.JoinVertical(lipgloss.Left, options, lipgloss.NewStyle().MarginLeft(lipgloss.Width(label)).Render("..."))
				break
			}

			if match := m.filterMatches[i]; match >= 0 {
				first := o[0:match]
				middle := o[match : match+len(m.field.Value())]
				end := o[match+len(m.field.Value()):]

				o = lipgloss.JoinHorizontal(lipgloss.Left, first, lipgloss.NewStyle().Foreground(m.styler.ActiveColor().Value()).Render(middle), end)
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
