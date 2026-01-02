package chat

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/style"
)

func Test_NewParameterModal_WithDefaultValues(t *testing.T) {
	params := []model.Parameter{
		{Name: "FOO", DefaultValue: "default_foo", Type: model.ParameterTypeValue},
		{Name: "BAR", DefaultValue: "default_bar", Type: model.ParameterTypeValue},
	}

	modal := NewParameterModal(params, style.Style{}, afero.NewMemMapFs())
	values := modal.GetValues()

	assert.Len(t, values, 2)
	assert.Equal(t, "default_foo", values[0])
	assert.Equal(t, "default_bar", values[1])
}

func Test_parameterModal_Update_Escape(t *testing.T) {
	params := []model.Parameter{
		{Name: "FOO", Type: model.ParameterTypeValue},
	}
	modal := NewParameterModal(params, style.Style{}, afero.NewMemMapFs())
	modal.Init()

	msg := tea.KeyMsg{Type: tea.KeyEscape}
	modal, _ = modal.Update(msg)

	assert.True(t, modal.IsCanceled())
	assert.False(t, modal.IsSubmitted())
}

func Test_parameterModal_Update_CtrlE_Submit(t *testing.T) {
	params := []model.Parameter{
		{Name: "FOO", Type: model.ParameterTypeValue},
	}
	modal := NewParameterModal(params, style.Style{}, afero.NewMemMapFs())
	modal.Init()

	msg := tea.KeyMsg{Type: tea.KeyCtrlE}
	modal, _ = modal.Update(msg)

	assert.True(t, modal.IsSubmitted())
	assert.False(t, modal.IsCanceled())
}

func Test_parameterModal_Update_Tab_Navigation(t *testing.T) {
	params := []model.Parameter{
		{Name: "FOO", Type: model.ParameterTypeValue},
		{Name: "BAR", Type: model.ParameterTypeValue},
	}
	modal := NewParameterModal(params, style.Style{}, afero.NewMemMapFs())
	modal.Init()

	assert.Equal(t, 0, modal.elementFocus)
	assert.Equal(t, focusFields, modal.focusArea)

	// Tab to next field
	msg := tea.KeyMsg{Type: tea.KeyTab}
	modal, _ = modal.Update(msg)
	assert.Equal(t, 1, modal.elementFocus)

	// Tab to buttons
	modal, _ = modal.Update(msg)
	assert.Equal(t, focusButtons, modal.focusArea)
}

func Test_parameterModal_EnterInButtons_Execute(t *testing.T) {
	modal := NewParameterModal([]model.Parameter{}, style.Style{}, afero.NewMemMapFs())
	modal.Init() // This sets focus to buttons since no fields

	assert.Equal(t, focusButtons, modal.focusArea)
	assert.Equal(t, 0, modal.buttonFocus) // Execute button

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	modal, _ = modal.Update(msg)

	assert.True(t, modal.IsSubmitted())
}

func Test_parameterModal_EnterInButtons_Cancel(t *testing.T) {
	modal := NewParameterModal([]model.Parameter{}, style.Style{}, afero.NewMemMapFs())
	modal.Init()
	modal.buttonFocus = 1 // Cancel button

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	modal, _ = modal.Update(msg)

	assert.True(t, modal.IsCanceled())
}

func Test_parameterModal_UpArrowFromButtons_GoesToLastField(t *testing.T) {
	params := []model.Parameter{
		{Name: "FOO", Type: model.ParameterTypeValue},
		{Name: "BAR", Type: model.ParameterTypeValue},
	}
	modal := NewParameterModal(params, style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Navigate to buttons
	modal.focusArea = focusButtons
	modal.buttonFocus = 0

	// Press up arrow
	msg := tea.KeyMsg{Type: tea.KeyUp}
	modal, _ = modal.Update(msg)

	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 1, modal.elementFocus) // Last field (index 1 of 2 fields)
}

func Test_parameterModal_TabOnLastButton_CyclesToFirstField(t *testing.T) {
	params := []model.Parameter{
		{Name: "FOO", Type: model.ParameterTypeValue},
	}
	modal := NewParameterModal(params, style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Navigate to last button
	modal.focusArea = focusButtons
	modal.buttonFocus = 1 // Cancel button

	// Press tab
	msg := tea.KeyMsg{Type: tea.KeyTab}
	modal, _ = modal.Update(msg)

	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 0, modal.elementFocus)
}

func Test_parameterModal_CompleteNavigationCycle(t *testing.T) {
	params := []model.Parameter{
		{Name: "FOO", Type: model.ParameterTypeValue},
		{Name: "BAR", Type: model.ParameterTypeValue},
	}
	modal := NewParameterModal(params, style.Style{}, afero.NewMemMapFs())
	modal.Init()

	// Start at first field
	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 0, modal.elementFocus)

	// Tab through fields and buttons
	tab := tea.KeyMsg{Type: tea.KeyTab}
	modal, _ = modal.Update(tab) // To field 1
	assert.Equal(t, 1, modal.elementFocus)

	modal, _ = modal.Update(tab) // To Execute button
	assert.Equal(t, focusButtons, modal.focusArea)
	assert.Equal(t, 0, modal.buttonFocus)

	modal, _ = modal.Update(tab) // To Cancel button
	assert.Equal(t, 1, modal.buttonFocus)

	modal, _ = modal.Update(tab) // Cycle back to first field
	assert.Equal(t, focusFields, modal.focusArea)
	assert.Equal(t, 0, modal.elementFocus)
}
