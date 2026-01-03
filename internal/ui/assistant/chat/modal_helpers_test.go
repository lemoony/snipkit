package chat

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func Test_handleModalKeyPress_Escape(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyEscape}
	action := handleModalKeyPress(msg, "e", focusFields, 0, 0, 2)
	assert.Equal(t, modalKeyCancel, action)
}

func Test_handleModalKeyPress_CtrlShortcut(t *testing.T) {
	tests := []struct {
		name     string
		shortcut string
		keyType  tea.KeyType
	}{
		{"ctrl+e", "e", tea.KeyCtrlE},
		{"ctrl+s", "s", tea.KeyCtrlS},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tt.keyType}
			action := handleModalKeyPress(msg, tt.shortcut, focusFields, 0, 0, 2)
			assert.Equal(t, modalKeySubmit, action)
		})
	}
}

func Test_handleModalKeyPress_EnterInButtons_Submit(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	action := handleModalKeyPress(msg, "e", focusButtons, 0, 0, 2)
	assert.Equal(t, modalKeySubmit, action)
}

func Test_handleModalKeyPress_EnterInButtons_Cancel(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	action := handleModalKeyPress(msg, "e", focusButtons, 1, 0, 2)
	assert.Equal(t, modalKeyCancel, action)
}

func Test_handleModalKeyPress_EnterInFields_LastField(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	action := handleModalKeyPress(msg, "e", focusFields, 0, 1, 2) // elementFocus=1 is last of 2 fields
	assert.Equal(t, modalKeyNavigateForward, action)
}

func Test_handleModalKeyPress_EnterInFields_NotLastField(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	action := handleModalKeyPress(msg, "e", focusFields, 0, 0, 2) // elementFocus=0 is not last of 2 fields
	assert.Equal(t, modalKeyNavigateFields, action)
}

func Test_handleModalKeyPress_Tab(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyTab}
	action := handleModalKeyPress(msg, "e", focusFields, 0, 0, 2)
	assert.Equal(t, modalKeyNavigateForward, action)
}

func Test_handleModalKeyPress_ShiftTab(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyShiftTab}
	action := handleModalKeyPress(msg, "e", focusFields, 0, 0, 2)
	assert.Equal(t, modalKeyNavigateBackward, action)
}

func Test_handleModalKeyPress_LeftInButtons(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyLeft}
	action := handleModalKeyPress(msg, "e", focusButtons, 1, 0, 2)
	assert.Equal(t, modalKeyNavigateButtonsBackward, action)
}

func Test_handleModalKeyPress_RightInButtons(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyRight}
	action := handleModalKeyPress(msg, "e", focusButtons, 0, 0, 2)
	assert.Equal(t, modalKeyNavigateButtonsForward, action)
}

func Test_handleModalKeyPress_LeftInFields(t *testing.T) {
	// Left/right in fields should delegate to field for text editing
	msg := tea.KeyMsg{Type: tea.KeyLeft}
	action := handleModalKeyPress(msg, "e", focusFields, 0, 0, 2)
	assert.Equal(t, modalKeyDelegateToField, action)
}

func Test_handleModalKeyPress_DefaultInFields(t *testing.T) {
	// Regular key press in fields should delegate
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	action := handleModalKeyPress(msg, "e", focusFields, 0, 0, 2)
	assert.Equal(t, modalKeyDelegateToField, action)
}
