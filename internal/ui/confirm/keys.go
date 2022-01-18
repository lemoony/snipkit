package confirm

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines keybindings. It satisfies to the help.KeyMap interface, which
// is used to render the menu menu.
type KeyMap struct {
	// Keybindings used when browsing the list.
	CursorUp   key.Binding
	CursorDown key.Binding
	Yes        key.Binding
	No         key.Binding
	Toggle     key.Binding
	Apply      key.Binding

	// The quit keybinding. This won't be caught when filtering.
	Quit key.Binding

	// The quit-no-matter-what keybinding. This will be caught when filtering.
	ForceQuit key.Binding
}

// defaultKeyMap returns a default set of keybindings.
func defaultKeyMap() KeyMap {
	return KeyMap{
		// Browsing.
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Yes: key.NewBinding(
			key.WithKeys("left", "y"),
			key.WithHelp("←/y", "yes"),
		),
		No: key.NewBinding(
			key.WithKeys("right", "n"),
			key.WithHelp("→/n", "no"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "toggle"),
		),
		Apply: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "apply"),
		),

		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// FullHelp returns bindings to show the full help view. It's part of the
// help.KeyMap interface.
func (m *model) FullHelp() [][]key.Binding {
	h := [][]key.Binding{{
		m.keyMap.Yes,
		m.keyMap.No,
		m.keyMap.Toggle,
		m.keyMap.Quit,
		m.keyMap.Apply,
	}}

	if m.isScrollable() {
		h = append([][]key.Binding{
			{
				m.keyMap.CursorUp,
				m.keyMap.CursorDown,
			},
		}, h...)
	}

	return h
}

// ShortHelp returns bindings to show in the abbreviated help view. It's part
// of the help.KeyMap interface.
func (m *model) ShortHelp() []key.Binding {
	h := []key.Binding{
		m.keyMap.Yes,
		m.keyMap.No,
		m.keyMap.Toggle,
		m.keyMap.Quit,
		m.keyMap.Apply,
	}

	if m.isScrollable() {
		h = append([]key.Binding{
			m.keyMap.CursorUp,
			m.keyMap.CursorDown,
		}, h...)
	}

	return h
}
