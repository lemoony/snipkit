package prompt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	History []string
}

type model struct {
	form    *huh.Form
	history []string

	quitting     bool
	success      bool
	latestPrompt string
}

func ShowPrompt(config Config, teaOptions ...tea.ProgramOption) (bool, string) {
	m := newModel(config, teaOptions...)

	teaModel, err := tea.NewProgram(m, teaOptions...).Run()
	if err != nil {
		return false, ""
	}

	resultModel := teaModel.(*model)
	return m.success, resultModel.latestPrompt
}

func newModel(config Config, teaOptions ...tea.ProgramOption) *model {
	m := &model{
		history: config.History,
		success: true,
	}

	placeholder := "Type here..."
	description := "What do you want the script to do?"
	if len(config.History) > 0 {
		var sb strings.Builder
		for i, v := range config.History {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(fmt.Sprintf("%s%s", lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Italic(true).Render(fmt.Sprintf("[%d] ", i+1)), v))
		}
		placeholder = "Type a new prompt or just press enter ..."
		description = fmt.Sprintf(
			"%s\n%s\n%s",
			"Do you want to provide additional context or change anything?",
			lipgloss.NewStyle().Italic(true).Render("Your previous prompts and their results are automatically provided as context:"),
			sb.String(),
		)
	}

	inputPrompt := huh.NewInput().
		Title("SnipKit Assistant").
		Description(description).
		Prompt("> ").
		Placeholder(placeholder).
		Value(&m.latestPrompt)

	m.form = huh.NewForm(huh.NewGroup(inputPrompt)).
		WithProgramOptions(teaOptions...).
		WithShowHelp(false)

	return m
}

func (m *model) Init() tea.Cmd {
	return m.form.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyCtrlC || keyMsg.Type == tea.KeyEsc {
			m.quitting = true
			m.success = false
			return m, tea.Quit
		}
	}

	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		m.quitting = true
		return m, tea.Quit
	}

	return m, cmd
}

func (m *model) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n%s", m.form.View())
}
