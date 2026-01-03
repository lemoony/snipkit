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
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("", "down"),
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
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "quit"),
		),
	}
}

// FieldKeyMap defines keybindings. It satisfies to the help.KeyMap interface, which
// is used to render the menu menu.
type FieldKeyMap struct {
	// Keybindings used when browsing the list.
	CursorUp        key.Binding
	CursorDown      key.Binding
	Apply           key.Binding
	ApplyCompletion key.Binding // Right arrow to apply without navigating

	// The quit keybinding. This won't be caught when filtering.
	Quit key.Binding
}

// defaultFieldKeyMap returns a default set of keybindings.
func defaultFieldKeyMap() FieldKeyMap {
	return FieldKeyMap{
		// Browsing.
		CursorUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Apply: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "apply"),
		),
		ApplyCompletion: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "complete"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "quit"),
		),
	}
}
