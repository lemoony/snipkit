package chat

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/style"
)

// Mode initialization tests - verify correct mode is set based on config

func Test_newUnifiedChatModel_EmptyConfig(t *testing.T) {
	config := UnifiedConfig{
		History:    []HistoryEntry{},
		Generating: false,
	}

	m := newUnifiedChatModel(config, style.Style{})

	assert.Equal(t, UIModeInput, m.currentMode)
	assert.Empty(t, m.messages)
	assert.False(t, m.generating)
}

func Test_newUnifiedChatModel_GeneratingMode(t *testing.T) {
	scriptChan := make(chan assistant.ParsedScript, 1)
	config := UnifiedConfig{
		History:    []HistoryEntry{{UserPrompt: "test prompt"}},
		Generating: true,
		ScriptChan: scriptChan,
	}

	m := newUnifiedChatModel(config, style.Style{})

	assert.Equal(t, UIModeGenerating, m.currentMode)
	assert.True(t, m.generating)
	// Should have user message + generating placeholder
	assert.Len(t, m.messages, 2)
}

func Test_newUnifiedChatModel_ScriptReadyMode(t *testing.T) {
	config := UnifiedConfig{
		History: []HistoryEntry{
			{UserPrompt: "test prompt", GeneratedScript: "echo hello"},
		},
	}

	m := newUnifiedChatModel(config, style.Style{})

	assert.Equal(t, UIModeScriptReady, m.currentMode)
	assert.NotNil(t, m.generatedScript)
}

func Test_newUnifiedChatModel_PostExecutionMode(t *testing.T) {
	exitCode := 0
	duration := 100 * time.Millisecond
	execTime := time.Now()
	config := UnifiedConfig{
		History: []HistoryEntry{
			{
				UserPrompt:      "test prompt",
				GeneratedScript: "echo hello",
				ExecutionOutput: "hello",
				ExitCode:        &exitCode,
				Duration:        &duration,
				ExecutionTime:   &execTime,
			},
		},
	}

	m := newUnifiedChatModel(config, style.Style{})

	assert.Equal(t, UIModePostExecution, m.currentMode)
	assert.True(t, m.hasExecutionOutput)
}

// State transition test

func Test_unifiedChatModel_TransitionToMode(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})
	m.ready = true
	m.width = 80
	m.height = 24

	m.transitionToMode(UIModeScriptReady)
	assert.Equal(t, UIModeScriptReady, m.currentMode)

	m.transitionToMode(UIModePostExecution)
	assert.Equal(t, UIModePostExecution, m.currentMode)
	assert.True(t, m.hasExecutionOutput)

	m.transitionToMode(UIModeInput)
	assert.Equal(t, UIModeInput, m.currentMode)
}

// Input mode tests - key handling

func Test_unifiedChatModel_HandleInputMode_Escape(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})
	m.setupInput()

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	_, _ = m.handleInputMode(msg)

	assert.True(t, m.quitting)
	assert.Equal(t, PreviewActionCancel, m.action)
}

func Test_unifiedChatModel_HandleInputMode_EnterWithText(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})
	m.setupInput()
	m.input.SetValue("test prompt")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _ = m.handleInputMode(msg)

	assert.True(t, m.quitting)
	assert.Equal(t, PreviewActionRevise, m.action)
	assert.Equal(t, "test prompt", m.latestPrompt)
}

func Test_unifiedChatModel_HandleInputMode_EnterEmpty(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})
	m.setupInput()

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := m.handleInputMode(msg)

	assert.False(t, m.quitting)
	assert.Nil(t, cmd)
}

// Generating mode test

func Test_unifiedChatModel_HandleGeneratingMode_Escape(t *testing.T) {
	scriptChan := make(chan assistant.ParsedScript, 1)
	config := UnifiedConfig{
		History:    []HistoryEntry{{UserPrompt: "test"}},
		Generating: true,
		ScriptChan: scriptChan,
	}
	m := newUnifiedChatModel(config, style.Style{})

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	_, _ = m.handleGeneratingMode(msg)

	assert.True(t, m.quitting)
	assert.Equal(t, PreviewActionCancel, m.action)
}

// Action bar tests - navigation and shortcuts

func Test_unifiedChatModel_HandleActionBarInput_Escape(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})

	options := m.getScriptReadyOptions()
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	shouldExecute, action, _ := m.handleActionBarInput(msg, options)

	assert.True(t, shouldExecute)
	assert.Equal(t, PreviewActionCancel, action)
}

func Test_unifiedChatModel_HandleActionBarInput_LeftRight(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})
	m.selectedOption = 0

	options := m.getScriptReadyOptions()

	msg := tea.KeyMsg{Type: tea.KeyRight}
	shouldExecute, _, newSelected := m.handleActionBarInput(msg, options)
	assert.False(t, shouldExecute)
	assert.Equal(t, 1, newSelected)

	m.selectedOption = 1
	msg = tea.KeyMsg{Type: tea.KeyLeft}
	shouldExecute, _, newSelected = m.handleActionBarInput(msg, options)
	assert.False(t, shouldExecute)
	assert.Equal(t, 0, newSelected)
}

func Test_unifiedChatModel_HandleActionBarInput_Enter(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})
	m.selectedOption = 0

	options := m.getScriptReadyOptions()
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	shouldExecute, action, _ := m.handleActionBarInput(msg, options)

	assert.True(t, shouldExecute)
	assert.Equal(t, PreviewActionExecute, action)
}

func Test_unifiedChatModel_HandleActionBarInput_Shortcuts(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})
	options := m.getScriptReadyOptions()

	tests := []struct {
		key      string
		expected PreviewAction
	}{
		{"e", PreviewActionExecute},
		{"o", PreviewActionEdit},
		{"r", PreviewActionRevise},
		{"c", PreviewActionCancel},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			shouldExecute, action, _ := m.handleActionBarInput(msg, options)
			assert.True(t, shouldExecute)
			assert.Equal(t, tt.expected, action)
		})
	}
}

// Script ready mode tests - user actions

func Test_unifiedChatModel_HandleScriptReadyMode_Execute(t *testing.T) {
	config := UnifiedConfig{
		History: []HistoryEntry{
			{UserPrompt: "test", GeneratedScript: "echo hi"},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})
	m.selectedOption = 0
	m.generatedScript = assistant.ParsedScript{Contents: "echo hi"}

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _ = m.handleScriptReadyMode(msg)

	assert.True(t, m.quitting)
	assert.Equal(t, PreviewActionExecute, m.action)
}

func Test_unifiedChatModel_HandleScriptReadyMode_Edit(t *testing.T) {
	config := UnifiedConfig{
		History: []HistoryEntry{
			{UserPrompt: "test", GeneratedScript: "echo hi"},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("o")}
	_, _ = m.handleScriptReadyMode(msg)

	assert.True(t, m.quitting)
	assert.Equal(t, PreviewActionEdit, m.action)
}

func Test_unifiedChatModel_HandleScriptReadyMode_Revise(t *testing.T) {
	config := UnifiedConfig{
		History: []HistoryEntry{
			{UserPrompt: "test", GeneratedScript: "echo hi"},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})
	m.ready = true
	m.width = 80
	m.height = 24
	m.handleWindowSize(tea.WindowSizeMsg{Width: 80, Height: 24})

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}
	_, _ = m.handleScriptReadyMode(msg)

	assert.Equal(t, UIModeInput, m.currentMode)
}

// Post-execution mode tests

func Test_unifiedChatModel_HandlePostExecutionMode_ExecuteAgain(t *testing.T) {
	exitCode := 0
	duration := 100 * time.Millisecond
	execTime := time.Now()
	config := UnifiedConfig{
		History: []HistoryEntry{
			{
				UserPrompt:      "test",
				GeneratedScript: "echo hi",
				ExecutionOutput: "hi",
				ExitCode:        &exitCode,
				Duration:        &duration,
				ExecutionTime:   &execTime,
			},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})
	m.selectedOption = 0

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, _ = m.handlePostExecutionMode(msg)

	assert.True(t, m.quitting)
	assert.Equal(t, PreviewActionExecute, m.action)
}

func Test_unifiedChatModel_HandlePostExecutionMode_ExitNoSave(t *testing.T) {
	exitCode := 0
	duration := 100 * time.Millisecond
	execTime := time.Now()
	config := UnifiedConfig{
		History: []HistoryEntry{
			{
				UserPrompt:      "test",
				GeneratedScript: "echo hi",
				ExecutionOutput: "hi",
				ExitCode:        &exitCode,
				Duration:        &duration,
				ExecutionTime:   &execTime,
			},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}
	_, _ = m.handlePostExecutionMode(msg)

	assert.True(t, m.quitting)
	assert.Equal(t, PreviewActionExitNoSave, m.action)
}

// Execute action tests - parameter modal triggering

func Test_unifiedChatModel_HandleExecuteAction_WithParameters(t *testing.T) {
	config := UnifiedConfig{
		History: []HistoryEntry{
			{UserPrompt: "test", GeneratedScript: "echo ${FOO}"},
		},
		Parameters: []model.Parameter{
			{Name: "FOO", Type: model.ParameterTypeValue},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})

	_, _ = m.handleExecuteAction()

	assert.Equal(t, modalParameters, m.modalState)
	assert.NotNil(t, m.paramModal)
}

func Test_unifiedChatModel_HandleExecuteAction_NoParameters(t *testing.T) {
	config := UnifiedConfig{
		History: []HistoryEntry{
			{UserPrompt: "test", GeneratedScript: "echo hi"},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})
	m.generatedScript = assistant.ParsedScript{Contents: "echo hi"}

	_, _ = m.handleExecuteAction()

	assert.True(t, m.quitting)
	assert.Equal(t, PreviewActionExecute, m.action)
}

// Window size handling tests

func Test_unifiedChatModel_HandleWindowSize_FirstSetup(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})

	assert.False(t, m.ready)

	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	m.handleWindowSize(msg)

	assert.True(t, m.ready)
	assert.Equal(t, 80, m.width)
	assert.Equal(t, 24, m.height)
}

func Test_unifiedChatModel_HandleWindowSize_Resize(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})

	// First setup
	m.handleWindowSize(tea.WindowSizeMsg{Width: 80, Height: 24})
	assert.True(t, m.ready)

	// Resize
	m.handleWindowSize(tea.WindowSizeMsg{Width: 120, Height: 40})

	assert.Equal(t, 120, m.width)
	assert.Equal(t, 40, m.height)
}

func Test_unifiedChatModel_HandleWindowSize_SmallTerminal(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})

	// Very small terminal
	m.handleWindowSize(tea.WindowSizeMsg{Width: 20, Height: 5})

	assert.True(t, m.ready)
	assert.Equal(t, 20, m.width)
	assert.Equal(t, 5, m.height)
}

func Test_unifiedChatModel_HandleWindowSize_InputMode(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})
	m.currentMode = UIModeInput
	m.setupInput()

	m.handleWindowSize(tea.WindowSizeMsg{Width: 80, Height: 24})

	assert.True(t, m.ready)
	assert.Greater(t, m.input.Width, 0)
}

func Test_unifiedChatModel_HandleWindowSize_GeneratingMode(t *testing.T) {
	scriptChan := make(chan assistant.ParsedScript, 1)
	config := UnifiedConfig{
		History:    []HistoryEntry{{UserPrompt: "test"}},
		Generating: true,
		ScriptChan: scriptChan,
	}
	m := newUnifiedChatModel(config, style.Style{})

	m.handleWindowSize(tea.WindowSizeMsg{Width: 80, Height: 24})

	assert.True(t, m.ready)
	assert.True(t, m.generating)
}

// Modal handling tests

func Test_unifiedChatModel_HandleModalUpdate_Parameters(t *testing.T) {
	config := UnifiedConfig{
		History: []HistoryEntry{
			{UserPrompt: "test", GeneratedScript: "echo ${FOO}"},
		},
		Parameters: []model.Parameter{
			{Name: "FOO", Type: model.ParameterTypeValue},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})
	m.handleWindowSize(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Trigger parameter modal
	m.handleExecuteAction()
	assert.Equal(t, modalParameters, m.modalState)

	// Update modal with key
	msg := tea.KeyMsg{Type: tea.KeyTab}
	m.handleModalUpdate(msg)

	// Modal should still be active (not submitted or canceled)
	assert.Equal(t, modalParameters, m.modalState)
}

func Test_unifiedChatModel_HandleModalUpdate_ParameterSubmit(t *testing.T) {
	config := UnifiedConfig{
		History: []HistoryEntry{
			{UserPrompt: "test", GeneratedScript: "echo ${FOO}"},
		},
		Parameters: []model.Parameter{
			{Name: "FOO", Type: model.ParameterTypeValue},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})
	m.handleWindowSize(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Trigger parameter modal
	m.handleExecuteAction()
	assert.Equal(t, modalParameters, m.modalState)
	assert.NotNil(t, m.paramModal)

	// Submit modal
	msg := tea.KeyMsg{Type: tea.KeyCtrlE}
	m.handleParameterModal(msg)

	assert.Equal(t, modalNone, m.modalState)
	assert.True(t, m.quitting)
	assert.Equal(t, PreviewActionExecute, m.action)
}

func Test_unifiedChatModel_HandleModalUpdate_ParameterCancel(t *testing.T) {
	config := UnifiedConfig{
		History: []HistoryEntry{
			{UserPrompt: "test", GeneratedScript: "echo ${FOO}"},
		},
		Parameters: []model.Parameter{
			{Name: "FOO", Type: model.ParameterTypeValue},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})
	m.handleWindowSize(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Trigger parameter modal
	m.handleExecuteAction()
	assert.Equal(t, modalParameters, m.modalState)

	// Cancel modal
	msg := tea.KeyMsg{Type: tea.KeyEscape}
	m.handleParameterModal(msg)

	assert.Equal(t, modalNone, m.modalState)
	assert.False(t, m.quitting)
}

func Test_unifiedChatModel_HandleModalUpdate_Executing(t *testing.T) {
	config := UnifiedConfig{
		History: []HistoryEntry{
			{UserPrompt: "test", GeneratedScript: "echo hi"},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})
	m.modalState = modalExecuting

	// Should return immediately
	_, cmd := m.handleModalUpdate(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Nil(t, cmd)
	assert.Equal(t, modalExecuting, m.modalState)
}

// Post-execution action tests

func Test_unifiedChatModel_HandlePostExecutionAction_SaveSnippet(t *testing.T) {
	exitCode := 0
	duration := 100 * time.Millisecond
	execTime := time.Now()
	config := UnifiedConfig{
		History: []HistoryEntry{
			{
				UserPrompt:      "test",
				GeneratedScript: "echo hi",
				ExecutionOutput: "hi",
				ExitCode:        &exitCode,
				Duration:        &duration,
				ExecutionTime:   &execTime,
			},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})
	m.handleWindowSize(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Press 's' to trigger save modal
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")}
	m.handlePostExecutionMode(msg)

	assert.Equal(t, modalSave, m.modalState)
	assert.NotNil(t, m.saveModal)
}

func Test_unifiedChatModel_HandlePostExecutionAction_Revise(t *testing.T) {
	exitCode := 0
	duration := 100 * time.Millisecond
	execTime := time.Now()
	config := UnifiedConfig{
		History: []HistoryEntry{
			{
				UserPrompt:      "test",
				GeneratedScript: "echo hi",
				ExecutionOutput: "hi",
				ExitCode:        &exitCode,
				Duration:        &duration,
				ExecutionTime:   &execTime,
			},
		},
	}
	m := newUnifiedChatModel(config, style.Style{})
	m.handleWindowSize(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Press 'r' to revise
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}
	m.handlePostExecutionMode(msg)

	assert.Equal(t, UIModeInput, m.currentMode)
}

// View rendering tests

func Test_unifiedChatModel_View_NotReady(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})

	view := m.View()

	// Should show loading/initializing message
	assert.NotEmpty(t, view)
}

func Test_unifiedChatModel_View_Ready(t *testing.T) {
	config := UnifiedConfig{History: []HistoryEntry{}}
	m := newUnifiedChatModel(config, style.Style{})
	m.handleWindowSize(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := m.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "SnipKit Assistant")
}
