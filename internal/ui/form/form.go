package form

import (
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	internalModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/form/field"
	"github.com/lemoony/snipkit/internal/ui/style"
)

const (
	buttonCount   = 2
	buttonPadding = 2
)

type model struct {
	colorProfile termenv.Profile

	input  *io.Reader
	output *io.Writer

	keyMap KeyMap
	help   help.Model

	width  int
	height int

	okButtonText string

	fields       []*field.Model
	elementFocus int
	showFields   int

	styler style.Style

	values []string
	apply  bool
}

func Show(parameters []internalModel.Parameter, okButton string, options ...Option) ([]string, bool) {
	m := initialModel(parameters, okButton)
	for _, o := range options {
		o.apply(m)
	}

	var teaOptions []tea.ProgramOption
	if m.input != nil {
		teaOptions = append(teaOptions, tea.WithInput(*m.input))
	}
	if m.output != nil {
		teaOptions = append(teaOptions, tea.WithOutput(*m.output))
	}

	teaOptions = append(teaOptions, tea.WithAltScreen())

	p := tea.NewProgram(m, teaOptions...)

	if err := p.Start(); err != nil {
		panic(err)
	}

	return m.values, m.apply
}

func initialModel(parameters []internalModel.Parameter, okButtonText string) *model {
	m := model{
		keyMap:       defaultKeyMap(),
		help:         help.New(),
		colorProfile: termenv.ColorProfile(),
		elementFocus: -1,
		showFields:   0,
		styler:       style.DefaultStyle(),
		okButtonText: okButtonText,
	}

	m.fields = make([]*field.Model, len(parameters))

	for i, f := range parameters {
		name := f.Key
		if f.Name != "" {
			name = f.Name
		}

		m.fields[i] = field.New(m.styler, name, f.Description, f.Values)
		if f.DefaultValue != "" {
			m.fields[i].SetValue(f.DefaultValue)
		}
	}

	for i := range m.fields {
		m.fields[i].SetLabelWidth(m.maxLabelWidth())
	}

	m.changeFocus()

	return &m
}

func (m *model) setResultValues() {
	m.values = make([]string, len(m.fields))
	for i := range m.fields {
		m.values[i] = m.fields[i].Value()
	}
}

func (m *model) changeFocus() tea.Cmd {
	nextFocus := m.elementFocus + 1
	if nextFocus >= len(m.fields)+buttonCount {
		nextFocus = 0
	}

	if m.showFields < nextFocus {
		m.showFields = nextFocus
	}

	if prev := m.elementFocus; prev >= 0 && prev < len(m.fields) {
		m.fields[prev].Blur()
	}

	m.elementFocus = nextFocus
	if nextFocus < len(m.fields) {
		return m.fields[nextFocus].Focus()
	}
	return nil
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Apply):
			switch {
			case m.elementFocus == len(m.fields):
				m.setResultValues()
				m.apply = true
				return m, tea.Quit
			case m.elementFocus == len(m.fields)+1:
				m.apply = false
				return m, tea.Quit
			default:
				cmds = append(cmds, m.changeFocus())
			}
		case key.Matches(msg, m.keyMap.Next):
			cmds = append(cmds, m.changeFocus())
		case msg.Type == tea.KeyRunes:
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.styler.SetSize(msg.Width, msg.Height)
	}

	if m.elementFocus < len(m.fields) {
		m.fields[m.elementFocus], cmd = m.fields[m.elementFocus].Update(msg)
	}

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) content() string {
	var sections []string
	sections = append(sections, m.styler.Title("This snippet requires parameters"))

	var fields []string
	for i, f := range m.fields {
		if i > m.showFields {
			break
		}
		fields = append(fields, m.styler.FormFieldWrapper(f.View()))
	}

	sections = append(sections, lipgloss.JoinVertical(lipgloss.Left, fields...))

	if m.showFields >= len(m.fields) {
		sections = append(sections, lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.renderButton(m.okButtonText, m.elementFocus == len(m.fields)),
			m.renderButton("Cancel", m.elementFocus == len(m.fields)+1),
		))
	}

	result := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return result
}

func (m *model) View() string {
	help := m.help.View(m)
	res := m.styler.MainView(m.content(), help, false)
	if m.styler.NeedsResize() {
		res = m.styler.MainView(m.content(), help, true)
	}
	return res
}

func (m *model) renderButton(text string, selected bool) string {
	return lipgloss.NewStyle().
		Margin(0, 1, 0, 0).
		Padding(0, buttonPadding).
		Foreground(m.styler.ButtonTextColor(selected)).
		Background(m.styler.ButtonColor(selected)).
		Render(text)
}

func (m *model) maxLabelWidth() int {
	result := 0
	for i := range m.fields {
		if w := lipgloss.Width(m.fields[i].Label); result < w {
			result = w
		}
	}
	return result
}
