package prompt

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lemoony/snipkit/internal/ui/style"
)

type Config struct {
	History []string
}

type model struct {
	history      []string
	input        textinput.Model
	quitting     bool
	success      bool
	latestPrompt string
	description  string

	styler style.Style
}

func ShowPrompt(config Config, styler style.Style, teaOptions ...tea.ProgramOption) (bool, string) {
	m := newModel(config, styler)

	if teaModel, err := tea.NewProgram(m, teaOptions...).Run(); err != nil {
		return false, ""
	} else if resultModel, ok := teaModel.(*model); ok {
		return resultModel.success, resultModel.latestPrompt
	}

	return false, ""
}

func newModel(config Config, styler style.Style) *model {
	m := &model{
		history: config.History,
		styler:  styler,
		success: true,
	}

	m.setupDescription()
	m.setupInput()

	return m
}

func (m *model) setupDescription() {
	if len(m.history) > 0 {
		m.description = fmt.Sprintf(
			"%s\n%s\n%s",
			m.styler.PromptLabel("Do you want to provide additional context or change anything?"),
			m.styler.PromptDescription("Your previous prompts and their results are automatically provided as context:"),
			m.renderHistory(),
		)
		m.input.Placeholder = "Type a new prompt or just press enter ..."
	} else {
		m.description = m.styler.PromptLabel("What do you want the script to do?")
	}
}

func (m *model) renderHistory() string {
	var sb strings.Builder
	for i, v := range m.history {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(m.styler.PromptDescription(fmt.Sprintf("%s%s", lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Italic(true).Render(fmt.Sprintf("[%d] ", i+1)), v)))
	}
	return sb.String()
}

func (m *model) setupInput() {
	m.input = textinput.New()
	m.input.Placeholder = "Type here..."
	m.input.Prompt = "> "
	m.input.Focus()
	m.input.PlaceholderStyle = lipgloss.NewStyle().Foreground(m.styler.PlaceholderColor().Value())
	m.input.PromptStyle = lipgloss.NewStyle().Foreground(m.styler.ActiveColor().Value())
	m.input.Cursor.Style = lipgloss.NewStyle().Foreground(m.styler.HighlightColor().Value())
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.input.Update(msg)
	m.input = form

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			m.success = false
			return m, tea.Quit
		case tea.KeyEnter:
			m.latestPrompt = m.input.Value()
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m *model) View() string {
	if m.quitting {
		return ""
	}

	return fmt.Sprintf(
		"\n%s\n%s\n%s",
		m.styler.Title("SnipKit Assistant"),
		m.description,
		m.styler.InputIndent(m.input.View()),
	)
}
