package spinner

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lemoony/snipkit/internal/ui/style"
)

type errMsg error

type stopMsg struct{}

type model struct {
	spinner  spinner.Model
	keyMap   KeyMap
	styler   style.Style
	quitting bool
	text     string
	title    string
	err      error
	stopChan chan bool
}

// KeyMap defines keybindings. It satisfies to the help.KeyMap interface, which
// is used to render the menu menu.
type KeyMap struct {
	// The quit keybinding. This won't be caught when filtering.
	Quit key.Binding
}

func initialModel(text string, title string, styler style.Style, stopChan chan bool) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		spinner:  s,
		text:     text,
		title:    title,
		styler:   styler,
		stopChan: stopChan,
		keyMap: KeyMap{
			Quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
	}
}

func (m *model) Init() tea.Cmd {
	// Start listening for stop signal
	return tea.Batch(m.spinner.Tick, waitForStop(m.stopChan))
}

func waitForStop(stopChan chan bool) tea.Cmd {
	return func() tea.Msg {
		<-stopChan // Wait for stop signal
		return stopMsg{}
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case stopMsg:
		// Handle stop signal
		m.quitting = true
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m *model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	var str string
	if !m.quitting {
		str = fmt.Sprintf("\n%s\n%s%s", m.styler.Title(m.title), m.spinner.View(), m.text)
	}
	return str
}

func ShowSpinner(text, title string, stopChan chan bool, styler style.Style, teaOptions ...tea.ProgramOption) {
	m := initialModel(text, title, styler, stopChan)
	p := tea.NewProgram(&m, teaOptions...)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
