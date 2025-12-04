package chat

import "time"

// HistoryEntry represents a complete interaction cycle in the assistant chat.
// It captures the user's prompt, the generated script, and the execution output.
type HistoryEntry struct {
	UserPrompt      string
	GeneratedScript string
	ExecutionOutput string         // Combined stdout + stderr
	ExitCode        *int           // Exit code from script execution (nil if not executed)
	Duration        *time.Duration // Execution duration (nil if not executed)
	ExecutionTime   *time.Time     // Timestamp when script was executed (nil if not executed)
}
