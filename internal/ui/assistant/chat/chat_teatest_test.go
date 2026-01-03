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

func Test_Chat_SubmitPrompt_Teatest(t *testing.T) {
	config := Config{History: []HistoryEntry{}}
	m := newModel(config, style.Style{})
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Wait for input prompt to appear
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte(">"))
	}, teatest.WithDuration(2*time.Second))

	// Type a prompt
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test prompt")})

	// Submit with Enter
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	finalModel := tm.FinalModel(t).(*chatModel)
	assert.True(t, finalModel.success)
	assert.Equal(t, "test prompt", finalModel.latestPrompt)
	assert.True(t, finalModel.quitting)
}

func Test_Chat_CancelPrompt_Teatest(t *testing.T) {
	config := Config{History: []HistoryEntry{}}
	m := newModel(config, style.Style{})
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Wait for input prompt to appear
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte(">"))
	}, teatest.WithDuration(2*time.Second))

	// Cancel with Escape
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	finalModel := tm.FinalModel(t).(*chatModel)
	assert.False(t, finalModel.success)
	assert.True(t, finalModel.quitting)
}

func Test_Chat_WithHistory_Teatest(t *testing.T) {
	// Create chat with existing history
	config := Config{
		History: []HistoryEntry{
			{UserPrompt: "previous prompt", GeneratedScript: "echo hello"},
		},
	}
	m := newModel(config, style.Style{})
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Wait for history to be rendered (should show the previous prompt)
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("previous prompt"))
	}, teatest.WithDuration(2*time.Second))

	// Type a new prompt and submit
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("new prompt")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	finalModel := tm.FinalModel(t).(*chatModel)
	assert.True(t, finalModel.success)
	assert.Equal(t, "new prompt", finalModel.latestPrompt)
}
