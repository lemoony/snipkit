package chat

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/ui/style"
)

// createScriptReadyConfig returns a config for a model in ScriptReady mode.
func createScriptReadyConfig() UnifiedConfig {
	return UnifiedConfig{
		History: []HistoryEntry{{
			UserPrompt:      "create echo script",
			GeneratedScript: "echo hello",
		}},
		Generating: false,
	}
}

func Test_UnifiedChat_ExecuteWorkflow_Teatest(t *testing.T) {
	m := newUnifiedChatModel(createScriptReadyConfig(), style.Style{})
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Wait for the action bar to appear
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Execute"))
	}, teatest.WithDuration(2*time.Second))

	// Press Enter to execute (Execute is default selected option)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	finalModel := tm.FinalModel(t).(*unifiedChatModel)
	assert.Equal(t, PreviewActionExecute, finalModel.action)
	assert.True(t, finalModel.quitting)
}

func Test_UnifiedChat_EditWorkflow_Teatest(t *testing.T) {
	m := newUnifiedChatModel(createScriptReadyConfig(), style.Style{})
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Wait for Open editor option
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Open editor"))
	}, teatest.WithDuration(2*time.Second))

	// Press 'o' for Open editor action (shortcut)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("o")})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	finalModel := tm.FinalModel(t).(*unifiedChatModel)
	assert.Equal(t, PreviewActionEdit, finalModel.action)
	assert.True(t, finalModel.quitting)
}

func Test_UnifiedChat_ReviseWorkflow_Teatest(t *testing.T) {
	m := newUnifiedChatModel(createScriptReadyConfig(), style.Style{})
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Wait for Revise option
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Revise"))
	}, teatest.WithDuration(2*time.Second))

	// Press 'r' for Revise action (shortcut) - transitions to input mode
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})

	// Wait for input mode (prompt indicator)
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte(">"))
	}, teatest.WithDuration(2*time.Second))

	// Type a new prompt and submit
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("make it faster")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	finalModel := tm.FinalModel(t).(*unifiedChatModel)
	assert.Equal(t, PreviewActionRevise, finalModel.action)
	assert.Equal(t, "make it faster", finalModel.latestPrompt)
	assert.True(t, finalModel.quitting)
}

func Test_UnifiedChat_CancelWorkflow_Teatest(t *testing.T) {
	m := newUnifiedChatModel(createScriptReadyConfig(), style.Style{})
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Wait for the action bar to appear
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Execute"))
	}, teatest.WithDuration(2*time.Second))

	// Press Escape to cancel
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	finalModel := tm.FinalModel(t).(*unifiedChatModel)
	assert.Equal(t, PreviewActionCancel, finalModel.action)
	assert.True(t, finalModel.quitting)
}
