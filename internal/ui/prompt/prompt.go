package prompt

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func ShowPrompt(placeholder string, teaOptions ...tea.ProgramOption) (bool, string) {
	m := initialModel(placeholder)
	p := tea.NewProgram(&m, teaOptions...)
	if _, err := p.Run(); err != nil {
		return true, ""
	}
	return m.success, m.value
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

	err      error
	value    string
	quitting bool
	success  bool
}

func initialModel(placeholder string) model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Prompt = "? "
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
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			m.success = false
			return m, tea.Quit
		case tea.KeyEnter:
			m.value = m.textInput.Value()
			m.quitting = true
			m.success = true
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

func (m *model) View() string {
	if m.quitting {
		return m.textInput.View()
	}

	return m.textInput.View()
}
