package chat

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// PreviewAction represents an action taken in the preview/chat interface.
type PreviewAction int

const (
	PreviewActionCancel PreviewAction = iota
	PreviewActionExecute
	PreviewActionEdit
	PreviewActionRevise
	PreviewActionExitNoSave
)

// modalState represents the current state of the modal overlay.
type modalState int

const (
	modalNone modalState = iota
	modalParameters
	modalExecuting
	modalSave
)

// Spinner and timing constants.
const (
	spinnerTickMillis = 100
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Message types for bubbletea.
type tickMsg struct{}

// scriptReadyMsg is sent when async script generation completes.
type scriptReadyMsg struct {
	script interface{}
}

// tick returns a command that sends a tick message after a delay.
func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*spinnerTickMillis, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

// waitForScript waits for a script to be generated and sends a message when ready.
func waitForScript(scriptChan chan interface{}) tea.Cmd {
	return func() tea.Msg {
		script := <-scriptChan
		return scriptReadyMsg{script: script}
	}
}
