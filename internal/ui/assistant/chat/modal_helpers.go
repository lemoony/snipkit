package chat

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lemoony/snipkit/internal/ui/style"
)

// renderModalControlBar renders a two-button control bar with the given labels and shortcuts.
func renderModalControlBar(
	styler style.Style,
	focusArea focusArea,
	buttonFocus int,
	primaryLabel string,
	primaryShortcut string,
	secondaryLabel string,
	secondaryShortcut string,
) string {
	buttons := []string{}

	// Primary button (Execute/Save)
	primaryStyle := lipgloss.NewStyle().Padding(0, 1)
	if focusArea == focusButtons && buttonFocus == 0 {
		primaryStyle = primaryStyle.
			Bold(true).
			Foreground(styler.HighlightColor().Value()).
			Background(styler.ActiveColor().Value())
	} else {
		primaryStyle = primaryStyle.Foreground(styler.TextColor().Value())
	}
	buttons = append(buttons, primaryStyle.Render("["+primaryShortcut+"] "+primaryLabel))

	// Secondary button (Cancel)
	secondaryStyle := lipgloss.NewStyle().Padding(0, 1)
	if focusArea == focusButtons && buttonFocus == 1 {
		secondaryStyle = secondaryStyle.
			Bold(true).
			Foreground(styler.HighlightColor().Value()).
			Background(styler.ActiveColor().Value())
	} else {
		secondaryStyle = secondaryStyle.Foreground(styler.TextColor().Value())
	}
	buttons = append(buttons, secondaryStyle.Render("["+secondaryShortcut+"] "+secondaryLabel))

	// Join buttons horizontally
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)

	// Add help text below
	helpText := lipgloss.NewStyle().
		Foreground(styler.PlaceholderColor().Value()).
		Render("Tab: navigate • Enter: select • " + primaryShortcut + "/" + secondaryShortcut + ": shortcuts")

	return lipgloss.JoinVertical(lipgloss.Left, buttonRow, helpText)
}

// modalKeyAction represents the result of handling a key in a modal.
type modalKeyAction int

const (
	modalKeyNone modalKeyAction = iota
	modalKeySubmit
	modalKeyCancel
	modalKeyNavigateForward
	modalKeyNavigateBackward
	modalKeyNavigateFields
	modalKeyNavigateButtonsForward
	modalKeyNavigateButtonsBackward
	modalKeyDelegateToField
)

// handleModalKeyPress processes common modal key presses and returns the appropriate action.
func handleModalKeyPress(
	keyMsg tea.KeyMsg,
	primaryShortcut string,
	focusArea focusArea,
	buttonFocus int,
	elementFocus int,
	fieldCount int,
) modalKeyAction {
	switch keyMsg.String() {
	case "esc", "c":
		return modalKeyCancel

	case primaryShortcut:
		return modalKeySubmit

	case "enter":
		if focusArea == focusButtons {
			if buttonFocus == 0 {
				return modalKeySubmit
			}
			return modalKeyCancel
		}
		// In field - check if it's the last field
		if elementFocus == fieldCount-1 {
			return modalKeyNavigateForward
		}
		return modalKeyNavigateFields

	case "tab", "down":
		return modalKeyNavigateForward

	case "shift+tab", "up":
		return modalKeyNavigateBackward

	case keyLeft, keyRight:
		if focusArea == focusButtons {
			if keyMsg.String() == keyLeft {
				return modalKeyNavigateButtonsBackward
			}
			return modalKeyNavigateButtonsForward
		}
	}

	// Delegate to field if in field area
	if focusArea == focusFields && elementFocus < fieldCount {
		return modalKeyDelegateToField
	}

	return modalKeyNone
}

// setupTextInput configures a text input with common styling and placeholder logic.
func setupTextInput(styler style.Style, hasMessages bool) textinput.Model {
	input := textinput.New()

	if hasMessages {
		input.Placeholder = "Type a new prompt or just press enter..."
	} else {
		input.Placeholder = "What do you want the script to do?"
	}

	input.Prompt = "> "
	input.Focus()
	input.PlaceholderStyle = lipgloss.NewStyle().
		Foreground(styler.PlaceholderColor().Value()).
		Background(styler.VerySubduedColor().Value()).
		Italic(true)
	input.PromptStyle = lipgloss.NewStyle().
		Foreground(styler.ActiveColor().Value()).
		Background(styler.VerySubduedColor().Value()).
		Bold(true)
	input.TextStyle = lipgloss.NewStyle().
		Foreground(styler.TextColor().Value()).
		Background(styler.VerySubduedColor().Value()).
		Bold(true)
	input.Cursor.Style = lipgloss.NewStyle().Foreground(styler.HighlightColor().Value())

	return input
}
