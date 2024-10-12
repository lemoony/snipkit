package prompt

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func ShowPrompt() (bool, string) {
	model := initialModel()
	p := tea.NewProgram(&model)
	if _, err := p.Run(); err != nil {
		return false, ""
	}
	return true, model.value
}

type (
	errMsg error
)

// KeyMap defines keybindings. It satisfies to the help.KeyMap interface, which
// is used to render the menu menu.
type KeyMap struct {
	// The quit keybinding. This won't be caught when filtering.
	Quit key.Binding
}

type model struct {
	textInput textinput.Model
	keyMap    KeyMap
	help      help.Model

	err      error
	value    string
	quitting bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "What do you want the script to do?"
	ti.Focus()

	return model{
		textInput: ti,
		err:       nil,
		keyMap: KeyMap{
			Quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
		help: help.New(),
	}
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			m.value = m.textInput.Value()
			m.quitting = true
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// FullHelp returns bindings to show the full help view. It's part of the
// help.KeyMap interface.
func (m *model) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

// ShortHelp returns bindings to show in the abbreviated help view. It's part
// of the help.KeyMap interface.
func (m *model) ShortHelp() []key.Binding {
	h := []key.Binding{
		m.keyMap.Quit,
	}
	return h
}

func (m *model) View() string {
	if m.quitting {
		return m.textInput.View()
	}

	return fmt.Sprintf(
		"%s\n\n%s",
		m.textInput.View(),
		m.help.View(m),
	) + "\n"
}
