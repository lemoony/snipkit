package chat

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/ui/style"
)

func Test_NewSaveModal_Initialization(t *testing.T) {
	modal := NewSaveModal("script.sh", "My Snippet", style.Style{}, afero.NewMemMapFs())

	assert.Equal(t, "script.sh", modal.GetFilename())
	assert.Equal(t, "My Snippet", modal.GetSnippetName())
	assert.False(t, modal.IsSubmitted())
	assert.False(t, modal.IsCanceled())
}

func Test_NewSaveModal_EmptyValues(t *testing.T) {
	modal := NewSaveModal("", "", style.Style{}, afero.NewMemMapFs())

	assert.Equal(t, "", modal.GetFilename())
	assert.Equal(t, "", modal.GetSnippetName())
}

func Test_saveModal_Init_FocusFirstField(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 0, modal.elementFocus)
	assert.Equal(t, 0, modal.buttonFocus)
}

func Test_saveModal_Update_Escape_SetsCancel(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	msg := tea.KeyMsg{Type: tea.KeyEscape}
	modal, _ = modal.Update(msg)

	assert.True(t, modal.IsCanceled())
	assert.False(t, modal.IsSubmitted())
}

func Test_saveModal_Update_CtrlS_SetsSubmit(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	msg := tea.KeyMsg{Type: tea.KeyCtrlS}
	modal, _ = modal.Update(msg)

	assert.True(t, modal.IsSubmitted())
	assert.False(t, modal.IsCanceled())
}

func Test_saveModal_Update_EnterOnSaveButton_SetsSubmit(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Navigate to buttons area with Save button focused
	modal.focusArea = focusButtons
	modal.buttonFocus = 0 // Save button

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	modal, _ = modal.Update(msg)

	assert.True(t, modal.IsSubmitted())
	assert.False(t, modal.IsCanceled())
}

func Test_saveModal_Update_EnterOnCancelButton_SetsCancel(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Navigate to buttons area with Cancel button focused
	modal.focusArea = focusButtons
	modal.buttonFocus = 1 // Cancel button

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	modal, _ = modal.Update(msg)

	assert.True(t, modal.IsCanceled())
	assert.False(t, modal.IsSubmitted())
}

func Test_saveModal_NavigateForward_ThroughFields(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Start at first field
	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 0, modal.elementFocus)

	// Tab to next field
	tab := tea.KeyMsg{Type: tea.KeyTab}
	modal, _ = modal.Update(tab)

	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 1, modal.elementFocus) // Second field
}

func Test_saveModal_NavigateForward_FieldToButtons(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Navigate to last field
	modal.elementFocus = 1
	modal.fields[0].Blur()
	modal.fields[1].Focus()

	// Tab to buttons
	tab := tea.KeyMsg{Type: tea.KeyTab}
	modal, _ = modal.Update(tab)

	assert.Equal(t, focusButtons, modal.focusArea)
	assert.Equal(t, 0, modal.buttonFocus) // Save button
}

func Test_saveModal_NavigateForward_CycleButtons(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	modal.focusArea = focusButtons
	modal.buttonFocus = 0 // Save button

	tab := tea.KeyMsg{Type: tea.KeyTab}
	modal, _ = modal.Update(tab)

	assert.Equal(t, 1, modal.buttonFocus) // Cancel button

	modal, _ = modal.Update(tab)

	// Should wrap back to first field
	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 0, modal.elementFocus)
}

func Test_saveModal_NavigateBackward_FromSecondField(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Start at second field
	modal.elementFocus = 1
	modal.fields[0].Blur()
	modal.fields[1].Focus()

	// Shift+Tab to previous field
	shiftTab := tea.KeyMsg{Type: tea.KeyShiftTab}
	modal, _ = modal.Update(shiftTab)

	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 0, modal.elementFocus)
}

func Test_saveModal_NavigateBackward_FromFirstField_WrapsToCancel(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Start at first field
	assert.Equal(t, 0, modal.elementFocus)

	// Shift+Tab should wrap to Cancel button
	shiftTab := tea.KeyMsg{Type: tea.KeyShiftTab}
	modal, _ = modal.Update(shiftTab)

	assert.Equal(t, focusButtons, modal.focusArea)
	assert.Equal(t, 1, modal.buttonFocus) // Cancel button
}

func Test_saveModal_NavigateBackward_FromButtons_ToLastField(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	modal.focusArea = focusButtons
	modal.buttonFocus = 0 // Save button

	// Shift+Tab should go to last field
	shiftTab := tea.KeyMsg{Type: tea.KeyShiftTab}
	modal, _ = modal.Update(shiftTab)

	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 1, modal.elementFocus) // Last field (index 1)
}

func Test_saveModal_ArrowKeys_CycleButtons(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	modal.focusArea = focusButtons
	modal.buttonFocus = 0

	// Right arrow to next button
	right := tea.KeyMsg{Type: tea.KeyRight}
	modal, _ = modal.Update(right)
	assert.Equal(t, 1, modal.buttonFocus)

	// Right arrow wraps
	modal, _ = modal.Update(right)
	assert.Equal(t, 0, modal.buttonFocus)

	// Left arrow to previous
	left := tea.KeyMsg{Type: tea.KeyLeft}
	modal, _ = modal.Update(left)
	assert.Equal(t, 1, modal.buttonFocus)
}

func Test_saveModal_UpArrow_FromButtons_GoesToLastField(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	modal.focusArea = focusButtons
	modal.buttonFocus = 0

	up := tea.KeyMsg{Type: tea.KeyUp}
	modal, _ = modal.Update(up)

	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 1, modal.elementFocus) // Last field
}

func Test_saveModal_GetFilename_ReturnsValue(t *testing.T) {
	modal := NewSaveModal("test-file.sh", "name", style.Style{}, afero.NewMemMapFs())

	assert.Equal(t, "test-file.sh", modal.GetFilename())
}

func Test_saveModal_GetSnippetName_ReturnsValue(t *testing.T) {
	modal := NewSaveModal("file.sh", "Test Snippet Name", style.Style{}, afero.NewMemMapFs())

	assert.Equal(t, "Test Snippet Name", modal.GetSnippetName())
}

func Test_saveModal_IsSubmitted_IsCanceled_InitialState(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())

	assert.False(t, modal.IsSubmitted())
	assert.False(t, modal.IsCanceled())
}

func Test_saveModal_CompleteNavigationCycle(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Start at first field
	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 0, modal.elementFocus)

	tab := tea.KeyMsg{Type: tea.KeyTab}

	// Tab through: field0 -> field1 -> Save -> Cancel -> field0
	modal, _ = modal.Update(tab) // To field 1
	assert.Equal(t, 1, modal.elementFocus)
	assert.Equal(t, focusFields, modal.focusArea)

	modal, _ = modal.Update(tab) // To Save button
	assert.Equal(t, focusButtons, modal.focusArea)
	assert.Equal(t, 0, modal.buttonFocus)

	modal, _ = modal.Update(tab) // To Cancel button
	assert.Equal(t, 1, modal.buttonFocus)

	modal, _ = modal.Update(tab) // Cycle back to first field
	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 0, modal.elementFocus)
}

func Test_saveModal_View_ReturnsNonEmpty(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	view := modal.View(80, 24)

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Save Snippet")
}

func Test_saveModal_NavigateFields_Forward(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Test internal navigateFields function
	modal.navigateFields(false) // Forward
	assert.Equal(t, 1, modal.elementFocus)

	modal.navigateFields(false) // Wraps to 0
	assert.Equal(t, 0, modal.elementFocus)
}

func Test_saveModal_NavigateFields_Backward(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Test backward navigation
	modal.navigateFields(true) // From 0 to last (wraps)
	assert.Equal(t, 1, modal.elementFocus)

	modal.navigateFields(true) // Back to 0
	assert.Equal(t, 0, modal.elementFocus)
}

func Test_saveModal_NavigateButtonsForward(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.focusArea = focusButtons
	modal.buttonFocus = 0

	modal.navigateButtonsForward()
	assert.Equal(t, 1, modal.buttonFocus)

	modal.navigateButtonsForward()
	assert.Equal(t, 0, modal.buttonFocus) // Wraps
}

func Test_saveModal_NavigateButtonsBackward(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.focusArea = focusButtons
	modal.buttonFocus = 0

	modal.navigateButtonsBackward()
	assert.Equal(t, 1, modal.buttonFocus) // Wraps

	modal.navigateButtonsBackward()
	assert.Equal(t, 0, modal.buttonFocus)
}

func Test_saveModal_NavigateUpFromButtons_ReturnsNil_WhenNotInButtons(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.Init()
	// Focus is on fields, not buttons

	cmd := modal.navigateUpFromButtons()
	assert.Nil(t, cmd)
}

func Test_saveModal_NavigateToFirstField(t *testing.T) {
	modal := NewSaveModal("file.sh", "name", style.Style{}, afero.NewMemMapFs())
	modal.focusArea = focusButtons
	modal.buttonFocus = 1

	modal.navigateToFirstField()

	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 0, modal.elementFocus)
}
