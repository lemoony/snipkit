package chat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_buildMessagesFromHistory_EmptyHistory(t *testing.T) {
	messages := buildMessagesFromHistory([]HistoryEntry{})
	assert.Empty(t, messages)
}

func Test_buildMessagesFromHistory_PromptOnly(t *testing.T) {
	history := []HistoryEntry{
		{UserPrompt: "test prompt"},
	}

	messages := buildMessagesFromHistory(history)

	assert.Len(t, messages, 1)
	assert.Equal(t, MessageTypeUser, messages[0].Type)
	assert.Equal(t, "test prompt", messages[0].Content)
}

func Test_buildMessagesFromHistory_PromptAndScript(t *testing.T) {
	history := []HistoryEntry{
		{
			UserPrompt:      "create a script",
			GeneratedScript: "echo hello",
		},
	}

	messages := buildMessagesFromHistory(history)

	assert.Len(t, messages, 2)
	assert.Equal(t, MessageTypeUser, messages[0].Type)
	assert.Equal(t, "create a script", messages[0].Content)
	assert.Equal(t, MessageTypeScript, messages[1].Type)
	assert.Equal(t, "echo hello", messages[1].Content)
}

func Test_buildMessagesFromHistory_FullCycle(t *testing.T) {
	exitCode := 0
	duration := 100 * time.Millisecond
	execTime := time.Now()

	history := []HistoryEntry{
		{
			UserPrompt:      "run something",
			GeneratedScript: "echo test",
			ExecutionOutput: "test output",
			ExitCode:        &exitCode,
			Duration:        &duration,
			ExecutionTime:   &execTime,
		},
	}

	messages := buildMessagesFromHistory(history)

	assert.Len(t, messages, 3)
	assert.Equal(t, MessageTypeUser, messages[0].Type)
	assert.Equal(t, MessageTypeScript, messages[1].Type)
	assert.Equal(t, MessageTypeOutput, messages[2].Type)
	assert.Equal(t, "test output", messages[2].Content)
	assert.Equal(t, &exitCode, messages[2].ExitCode)
	assert.Equal(t, &duration, messages[2].Duration)
}

func Test_buildMessagesFromHistory_SkipsDuplicateScripts(t *testing.T) {
	exitCode := 0
	duration := 50 * time.Millisecond
	execTime := time.Now()

	history := []HistoryEntry{
		{
			UserPrompt:      "first prompt",
			GeneratedScript: "echo hello",
			ExecutionOutput: "output1",
			ExitCode:        &exitCode,
			Duration:        &duration,
			ExecutionTime:   &execTime,
		},
		{
			UserPrompt:      "",           // Execute again (no new prompt)
			GeneratedScript: "echo hello", // Same script
			ExecutionOutput: "output2",
			ExitCode:        &exitCode,
			Duration:        &duration,
			ExecutionTime:   &execTime,
		},
	}

	messages := buildMessagesFromHistory(history)

	// Should have: prompt, script, output1, output2 (no duplicate script)
	assert.Len(t, messages, 4)
	assert.Equal(t, MessageTypeUser, messages[0].Type)
	assert.Equal(t, MessageTypeScript, messages[1].Type)
	assert.Equal(t, MessageTypeOutput, messages[2].Type)
	assert.Equal(t, "output1", messages[2].Content)
	assert.Equal(t, MessageTypeOutput, messages[3].Type)
	assert.Equal(t, "output2", messages[3].Content)
}

func Test_buildMessagesFromHistory_MultipleEntries(t *testing.T) {
	history := []HistoryEntry{
		{
			UserPrompt:      "first",
			GeneratedScript: "script1",
		},
		{
			UserPrompt:      "second",
			GeneratedScript: "script2",
		},
	}

	messages := buildMessagesFromHistory(history)

	assert.Len(t, messages, 4)
	assert.Equal(t, "first", messages[0].Content)
	assert.Equal(t, "script1", messages[1].Content)
	assert.Equal(t, "second", messages[2].Content)
	assert.Equal(t, "script2", messages[3].Content)
}

func Test_buildMessagesFromHistory_EmptyPromptSkipped(t *testing.T) {
	history := []HistoryEntry{
		{
			UserPrompt:      "", // Empty prompt
			GeneratedScript: "echo hello",
		},
	}

	messages := buildMessagesFromHistory(history)

	// Only script, no user message for empty prompt
	assert.Len(t, messages, 1)
	assert.Equal(t, MessageTypeScript, messages[0].Type)
}

func Test_buildMessagesFromHistory_EmptyOutput(t *testing.T) {
	exitCode := 0
	duration := 100 * time.Millisecond
	execTime := time.Now()

	tests := []struct {
		name   string
		output string
	}{
		{"empty string", ""},
		{"whitespace only", "   "},
		{"newlines only", "\n\n\n"},
		{"tabs and newlines", "\t\n\t\n"},
		{"mixed whitespace", "  \n\t  \n  "},
		{"ansi codes only", "\x1b[31m\x1b[0m"},
		{"ansi with whitespace", "\x1b[31m  \n  \x1b[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			history := []HistoryEntry{
				{
					UserPrompt:      "run something",
					GeneratedScript: "echo -n",
					ExecutionOutput: tt.output,
					ExitCode:        &exitCode,
					Duration:        &duration,
					ExecutionTime:   &execTime,
				},
			}

			messages := buildMessagesFromHistory(history)

			assert.Len(t, messages, 3)
			assert.Equal(t, MessageTypeUser, messages[0].Type)
			assert.Equal(t, MessageTypeScript, messages[1].Type)
			assert.Equal(t, MessageTypeOutput, messages[2].Type)
			assert.Equal(t, "Empty output", messages[2].Content)
			assert.Equal(t, &exitCode, messages[2].ExitCode)
			assert.Equal(t, &duration, messages[2].Duration)
		})
	}
}

func Test_buildMessagesFromHistory_FailedCommandEmptyOutput(t *testing.T) {
	failedExitCode := 1
	duration := 100 * time.Millisecond
	execTime := time.Now()

	tests := []struct {
		name   string
		output string
	}{
		{"empty string", ""},
		{"whitespace only", "   "},
		{"newlines only", "\n\n\n"},
		{"ansi codes only", "\x1b[31m\x1b[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			history := []HistoryEntry{
				{
					UserPrompt:      "run something",
					GeneratedScript: "false",
					ExecutionOutput: tt.output,
					ExitCode:        &failedExitCode,
					Duration:        &duration,
					ExecutionTime:   &execTime,
				},
			}

			messages := buildMessagesFromHistory(history)

			assert.Len(t, messages, 3)
			assert.Equal(t, MessageTypeUser, messages[0].Type)
			assert.Equal(t, MessageTypeScript, messages[1].Type)
			assert.Equal(t, MessageTypeOutput, messages[2].Type)
			// Failed commands should NOT show "Empty output"
			assert.Equal(t, "", messages[2].Content)
			assert.Equal(t, &failedExitCode, messages[2].ExitCode)
			assert.Equal(t, &duration, messages[2].Duration)
		})
	}
}

func Test_buildMessagesFromHistory_FailedCommandWithError(t *testing.T) {
	failedExitCode := 128
	duration := 100 * time.Millisecond
	execTime := time.Now()

	history := []HistoryEntry{
		{
			UserPrompt:      "git log on nonexistent file",
			GeneratedScript: "git log --oneline --follow -f nonexistent.txt",
			ExecutionOutput: "fatal: ambiguous argument 'nonexistent.txt': unknown revision or path not in the working tree.\nUse '--' to separate paths from revisions",
			ExitCode:        &failedExitCode,
			Duration:        &duration,
			ExecutionTime:   &execTime,
		},
	}

	messages := buildMessagesFromHistory(history)

	assert.Len(t, messages, 3)
	assert.Equal(t, MessageTypeUser, messages[0].Type)
	assert.Equal(t, MessageTypeScript, messages[1].Type)
	assert.Equal(t, MessageTypeOutput, messages[2].Type)
	// Failed commands should show the error message
	assert.Contains(t, messages[2].Content, "fatal: ambiguous argument")
	assert.Equal(t, &failedExitCode, messages[2].ExitCode)
	assert.Equal(t, &duration, messages[2].Duration)
}

func Test_stripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no ansi", "hello world", "hello world"},
		{"color code", "\x1b[31mred\x1b[0m", "red"},
		{"bold", "\x1b[1mbold\x1b[0m", "bold"},
		{"cursor movement", "\x1b[2Aup two lines", "up two lines"},
		{"clear line", "\x1b[2Kcleared", "cleared"},
		{"multiple codes", "\x1b[32m\x1b[1mgreen bold\x1b[0m", "green bold"},
		{"osc sequence", "\x1b]0;title\x07text", "text"},
		{"carriage return", "line1\rline2", "line1line2"},
		{"mixed content", "before\x1b[31mred\x1b[0mafter", "beforeredafter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripANSI(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_truncateContent_Truncated(t *testing.T) {
	content := "line1\nline2\nline3\nline4\nline5"
	result := truncateContent(content, 3)

	assert.Contains(t, result, "line1")
	assert.Contains(t, result, "line2")
	assert.Contains(t, result, "line3")
	assert.Contains(t, result, "... (2 more lines)")
	assert.NotContains(t, result, "line4")
	assert.NotContains(t, result, "line5")
}

func Test_truncateContent_ManyLines(t *testing.T) {
	var lines []string
	for i := 1; i <= 50; i++ {
		lines = append(lines, "line")
	}
	content := ""
	for i, l := range lines {
		if i > 0 {
			content += "\n"
		}
		content += l
	}

	result := truncateContent(content, 20)

	assert.Contains(t, result, "... (30 more lines)")
}
