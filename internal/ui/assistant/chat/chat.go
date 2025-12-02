package chat

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lemoony/snipkit/internal/ui/style"
)

type Config struct {
	History []HistoryEntry
}

type chatModel struct {
	messages []ChatMessage
	viewport viewport.Model
	input    textinput.Model

	width  int
	height int
	ready  bool

	quitting     bool
	success      bool
	latestPrompt string

	styler style.Style
}

func ShowChat(config Config, styler style.Style, teaOptions ...tea.ProgramOption) (bool, string) {
	m := newModel(config, styler)

	if teaModel, err := tea.NewProgram(m, teaOptions...).Run(); err != nil {
		return false, ""
	} else if resultModel, ok := teaModel.(*chatModel); ok {
		return resultModel.success, resultModel.latestPrompt
	}

	return false, ""
}

func newModel(config Config, styler style.Style) *chatModel {
	m := &chatModel{
		messages: buildMessagesFromHistory(config.History),
		styler:   styler,
		success:  true,
	}

	m.setupInput()

	return m
}

func (m *chatModel) setupInput() {
	m.input = textinput.New()

	if len(m.messages) > 0 {
		m.input.Placeholder = "Type a new prompt or just press enter..."
	} else {
		m.input.Placeholder = "What do you want the script to do?"
	}

	// Subtle grey background for input
	subtleGrey := lipgloss.Color("#2a2a2a")

	m.input.Prompt = "> "
	m.input.Focus()
	m.input.PlaceholderStyle = lipgloss.NewStyle().
		Foreground(m.styler.PlaceholderColor().Value()).
		Background(subtleGrey).
		Italic(true)
	m.input.PromptStyle = lipgloss.NewStyle().
		Foreground(m.styler.ActiveColor().Value()).
		Background(subtleGrey).
		Bold(true)
	m.input.TextStyle = lipgloss.NewStyle().
		Foreground(m.styler.TextColor().Value()).
		Background(subtleGrey).
		Bold(true)
	m.input.Cursor.Style = lipgloss.NewStyle().Foreground(m.styler.HighlightColor().Value())
}

func (m *chatModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			m.success = false
			return m, tea.Quit

		case tea.KeyEnter:
			m.latestPrompt = m.input.Value()
			m.quitting = true
			return m, tea.Quit

		case tea.KeyPgUp, tea.KeyPgDown:
			// Pass to viewport for scrolling
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd

		default:
			// Pass to input
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate dimensions
		titleHeight := lipgloss.Height(m.styler.Title("SnipKit Assistant"))
		inputHeight := 3 // input line + some spacing
		margins := 4     // top and bottom margins

		viewportHeight := msg.Height - titleHeight - inputHeight - margins
		if viewportHeight < 1 {
			viewportHeight = 1
		}

		if !m.ready {
			// First time setup
			m.viewport = viewport.New(msg.Width, viewportHeight)
			m.viewport.YPosition = 0
			m.viewport.SetContent(renderMessages(m.messages, m.styler, msg.Width))
			m.viewport.GotoBottom() // Start at bottom (latest messages)
			m.ready = true
		} else {
			// Resize existing viewport
			m.viewport.Width = msg.Width
			m.viewport.Height = viewportHeight
			m.viewport.SetContent(renderMessages(m.messages, m.styler, msg.Width))
		}

		// Update input width (ensure it's positive)
		inputWidth := msg.Width - len(m.input.Prompt) - 2
		if inputWidth < 1 {
			inputWidth = 1
		}
		m.input.Width = inputWidth

		return m, nil
	}

	return m, cmd
}

func (m *chatModel) View() string {
	if m.quitting {
		return ""
	}

	if !m.ready {
		return "\n  Initializing..."
	}

	// Build the view
	var sections []string

	// Title
	sections = append(sections, m.styler.Title("SnipKit Assistant"))

	// Viewport with message history
	sections = append(sections, m.viewport.View())

	// Help text if history is scrollable
	if m.viewport.TotalLineCount() > m.viewport.Height {
		helpText := lipgloss.NewStyle().
			Foreground(m.styler.PlaceholderColor().Value()).
			Render("  Use PgUp/PgDown to scroll • Esc to cancel • Enter to submit")
		sections = append(sections, helpText)
	}

	// Input area with enhanced styling
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(m.styler.ActiveColor().Value()).
		Background(lipgloss.Color("#2a2a2a")).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1).
		Width(m.width)

	sections = append(sections, inputStyle.Render(m.input.View()))

	return fmt.Sprintf("\n%s\n", lipgloss.JoinVertical(lipgloss.Left, sections...))
}
