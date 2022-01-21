package form

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
	Apply      key.Binding
	Next       key.Binding

	// The quit keybinding. This won't be caught when filtering.
	Quit key.Binding
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
		Apply: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "apply"),
		),
		Next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
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
	return [][]key.Binding{}
}

// ShortHelp returns bindings to show in the abbreviated help view. It's part
// of the help.KeyMap interface.
func (m *model) ShortHelp() []key.Binding {
	h := []key.Binding{
		m.keyMap.Next,
		m.keyMap.Quit,
		m.keyMap.Apply,
	}

	return h
}
