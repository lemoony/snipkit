package chat

import (
	"strings"

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
	buttons = append(buttons, primaryStyle.Render(primaryLabel))

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
	buttons = append(buttons, secondaryStyle.Render(secondaryLabel))

	// Join buttons horizontally
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)

	// Add help text below with Ctrl shortcuts
	helpText := lipgloss.NewStyle().
		Foreground(styler.PlaceholderColor().Value()).
		Render("Ctrl+" + strings.ToUpper(primaryShortcut) + ": " + strings.ToLower(primaryLabel) + " • Esc: cancel • Tab/↑/↓: navigate • Enter: select")

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
	modalKeyNavigateUp
	modalKeyNavigateToFirstField
)

// handleEnterKey processes the enter key based on focus area and position.
func handleEnterKey(focusArea focusArea, buttonFocus int, elementFocus int, fieldCount int) modalKeyAction {
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
}

// handleTabKey processes the tab key based on focus area and button position.
func handleTabKey(focusArea focusArea, buttonFocus int) modalKeyAction {
	if focusArea == focusButtons && buttonFocus == 1 {
		// Last button - cycle to first field
		return modalKeyNavigateToFirstField
	}
	return modalKeyNavigateForward
}

// handleArrowKeys processes arrow keys based on focus area and key pressed.
func handleArrowKeys(keyStr string, focusArea focusArea) modalKeyAction {
	if keyStr == "up" && focusArea == focusButtons {
		return modalKeyNavigateUp
	}
	if focusArea == focusButtons && (keyStr == keyLeft || keyStr == keyRight) {
		if keyStr == keyLeft {
			return modalKeyNavigateButtonsBackward
		}
		return modalKeyNavigateButtonsForward
	}
	return modalKeyNone
}

// handleModalKeyPress processes common modal key presses and returns the appropriate action.
func handleModalKeyPress(
	keyMsg tea.KeyMsg,
	primaryShortcut string,
	focusArea focusArea,
	buttonFocus int,
	elementFocus int,
	fieldCount int,
) modalKeyAction {
	keyStr := keyMsg.String()

	switch keyStr {
	case "esc":
		return modalKeyCancel
	case "ctrl+" + primaryShortcut:
		return modalKeySubmit
	case "enter":
		return handleEnterKey(focusArea, buttonFocus, elementFocus, fieldCount)
	case "tab":
		return handleTabKey(focusArea, buttonFocus)
	case "shift+tab":
		return modalKeyNavigateBackward
	case "up", keyLeft, keyRight:
		if action := handleArrowKeys(keyStr, focusArea); action != modalKeyNone {
			return action
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
