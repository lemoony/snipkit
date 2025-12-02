package chat

// HistoryEntry represents a complete interaction cycle in the assistant chat.
// It captures the user's prompt, the generated script, and the execution output.
type HistoryEntry struct {
	UserPrompt      string
	GeneratedScript string
	ExecutionOutput string // Combined stdout + stderr
}
