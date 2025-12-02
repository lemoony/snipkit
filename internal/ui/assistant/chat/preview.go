package chat

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lemoony/snipkit/internal/ui/style"
)

// PreviewConfig contains configuration for showing a script preview.
type PreviewConfig struct {
	History []HistoryEntry
	Script  string
}

// PreviewAction represents the user's choice after viewing the script preview.
type PreviewAction int

const (
	PreviewActionCancel PreviewAction = iota
	PreviewActionExecute
	PreviewActionEdit
	PreviewActionRevise
)

const (
	menuHeight = 3
)

var menuOptions = []struct {
	key    string
	label  string
	action PreviewAction
}{
	{"E", "Execute", PreviewActionExecute},
	{"O", "Open in editor", PreviewActionEdit},
	{"R", "Revise prompt", PreviewActionRevise},
	{"C", "Cancel", PreviewActionCancel},
}

type previewModel struct {
	messages []ChatMessage
	viewport viewport.Model

	width  int
	height int
	ready  bool

	quitting       bool
	action         PreviewAction
	selectedOption int // Currently selected menu option (0-3)

	styler style.Style
}

// ShowScriptPreview shows the conversation history with the generated script.
// Returns the user's chosen action: execute directly, open in editor, or cancel.
func ShowScriptPreview(config PreviewConfig, styler style.Style, teaOptions ...tea.ProgramOption) PreviewAction {
	m := newPreviewModel(config, styler)

	if teaModel, err := tea.NewProgram(m, teaOptions...).Run(); err != nil {
		return PreviewActionCancel
	} else if resultModel, ok := teaModel.(*previewModel); ok {
		return resultModel.action
	}

	return PreviewActionCancel
}

func newPreviewModel(config PreviewConfig, styler style.Style) *previewModel {
	messages := buildMessagesFromHistory(config.History)

	// Add the generated script as a message
	if config.Script != "" {
		messages = append(messages, ChatMessage{
			Type:    MessageTypeScript,
			Content: config.Script,
		})
	}

	return &previewModel{
		messages:       messages,
		styler:         styler,
		action:         PreviewActionCancel,
		selectedOption: 0, // Default to Execute (first menu item)
	}
}

func (m *previewModel) Init() tea.Cmd {
	return nil
}

func (m *previewModel) selectAction(action PreviewAction) (tea.Model, tea.Cmd) {
	m.quitting = true
	m.action = action
	return m, tea.Quit
}

func (m *previewModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m.selectAction(PreviewActionCancel)

	case tea.KeyEnter:
		return m.selectAction(menuOptions[m.selectedOption].action)

	case tea.KeyUp, tea.KeyLeft:
		m.selectedOption--
		if m.selectedOption < 0 {
			m.selectedOption = len(menuOptions) - 1
		}
		return m, nil

	case tea.KeyDown, tea.KeyRight:
		m.selectedOption++
		if m.selectedOption >= len(menuOptions) {
			m.selectedOption = 0
		}
		return m, nil

	case tea.KeyPgUp, tea.KeyPgDown:
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	// Keyboard shortcuts
	return m.handleShortcut(msg.String())
}

func (m *previewModel) handleShortcut(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "e", "E":
		return m.selectAction(PreviewActionExecute)
	case "r", "R":
		return m.selectAction(PreviewActionRevise)
	case "o", "O":
		return m.selectAction(PreviewActionEdit)
	case "c", "C":
		return m.selectAction(PreviewActionCancel)
	}
	return m, nil
}

func (m *previewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate dimensions
		titleHeight := lipgloss.Height(m.styler.Title("Generated Script"))
		margins := 4 // top and bottom margins

		viewportHeight := msg.Height - titleHeight - menuHeight - margins
		if viewportHeight < 1 {
			viewportHeight = 1
		}

		if !m.ready {
			// First time setup
			m.viewport = viewport.New(msg.Width, viewportHeight)
			m.viewport.YPosition = 0
			m.viewport.SetContent(renderMessages(m.messages, m.styler, msg.Width))
			m.viewport.GotoBottom() // Start at bottom (show the generated script)
			m.ready = true
		} else {
			// Resize existing viewport
			m.viewport.Width = msg.Width
			m.viewport.Height = viewportHeight
			m.viewport.SetContent(renderMessages(m.messages, m.styler, msg.Width))
		}

		return m, nil
	}

	return m, nil
}

func (m *previewModel) View() string {
	if m.quitting {
		return ""
	}

	if !m.ready {
		return "\n  Initializing..."
	}

	// Build the view
	var sections []string

	// Title
	sections = append(sections, m.styler.Title("Generated Script"))

	// Viewport with message history and script
	sections = append(sections, m.viewport.View())

	// Interactive menu
	sections = append(sections, m.renderMenu())

	return fmt.Sprintf("\n%s\n", lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (m *previewModel) renderMenu() string {
	var menuItems []string
	for i, opt := range menuOptions {
		var style lipgloss.Style
		if i == m.selectedOption {
			// Selected item - highlighted
			style = lipgloss.NewStyle().
				Bold(true).
				Foreground(m.styler.HighlightColor().Value()).
				Background(m.styler.ActiveColor().Value()).
				Padding(0, 1)
		} else {
			// Unselected item
			style = lipgloss.NewStyle().
				Foreground(m.styler.TextColor().Value()).
				Padding(0, 1)
		}

		item := fmt.Sprintf("[%s] %s", opt.key, opt.label)
		menuItems = append(menuItems, style.Render(item))
	}

	menu := lipgloss.JoinHorizontal(lipgloss.Top, menuItems...)

	helpText := lipgloss.NewStyle().
		Foreground(m.styler.PlaceholderColor().Value()).
		Render("\n  ↑/↓ to select • Enter to confirm • Esc to cancel • PgUp/PgDn to scroll")

	return fmt.Sprintf("\n  %s%s", menu, helpText)
}
