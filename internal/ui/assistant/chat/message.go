package chat

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/lipgloss"

	"github.com/lemoony/snipkit/internal/ui/style"
)

// ansiRegex matches ANSI escape sequences and other terminal control codes.
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\][^\x07]*\x07|\x1b[PX^_][^\x1b]*\x1b\\|\r`)

type MessageType int

const (
	MessageTypeUser MessageType = iota
	MessageTypeAssistant
	MessageTypeScript
	MessageTypeOutput
)

type ChatMessage struct {
	Type          MessageType
	Content       string
	Timestamp     time.Time
	ExitCode      *int           // Exit code for output messages
	Duration      *time.Duration // Execution duration for output messages
	ExecutionTime *time.Time     // Execution timestamp for output messages
}

const maxMessagesPerEntry = 3 // User prompt, script, and output

// buildMessagesFromHistory converts a history of interactions into chat messages.
func buildMessagesFromHistory(history []HistoryEntry) []ChatMessage {
	if len(history) == 0 {
		return []ChatMessage{}
	}

	messages := make([]ChatMessage, 0, len(history)*maxMessagesPerEntry)
	var previousScript string

	for _, entry := range history {
		// Add user prompt (only if not empty)
		if entry.UserPrompt != "" {
			messages = append(messages, ChatMessage{
				Type:      MessageTypeUser,
				Content:   entry.UserPrompt,
				Timestamp: time.Now(),
			})
		}

		// Add generated script (truncated) only if it's different from the previous one
		if entry.GeneratedScript != "" && entry.GeneratedScript != previousScript {
			script := truncateContent(entry.GeneratedScript, maxLines)
			messages = append(messages, ChatMessage{
				Type:      MessageTypeScript,
				Content:   script,
				Timestamp: time.Now(),
			})
			previousScript = entry.GeneratedScript
		}

		// Add execution output (truncated) - show "Empty output" if execution occurred but no output
		if entry.ExecutionOutput != "" || entry.ExitCode != nil {
			output := entry.ExecutionOutput
			strippedOutput := strings.TrimSpace(stripANSI(output))

			// Only show "Empty output" for successful commands with no output
			// Failed commands should show whatever output they have (even if empty after stripping)
			if strippedOutput == "" {
				if entry.ExitCode != nil && *entry.ExitCode != 0 {
					// Failed command with no visible output - keep empty, exit code indicator will show
					output = ""
				} else {
					// Successful command with no output
					output = "Empty output"
				}
			} else {
				output = truncateContent(output, maxLines)
			}

			messages = append(messages, ChatMessage{
				Type:          MessageTypeOutput,
				Content:       output,
				Timestamp:     time.Now(),
				ExitCode:      entry.ExitCode,
				Duration:      entry.Duration,
				ExecutionTime: entry.ExecutionTime,
			})
		}
	}

	return messages
}

// renderMessages renders all messages into a single string for the viewport.
func renderMessages(messages []ChatMessage, styler style.Style, width int) string {
	if len(messages) == 0 {
		contextNote := styler.PromptDescription(
			"The history and script results are automatically provided as context",
		)
		return contextNote
	}

	var sections []string
	for _, msg := range messages {
		sections = append(sections, renderMessage(msg, styler, width))
	}

	return strings.Join(sections, "\n\n")
}

// renderMessage renders a single message based on its type.
func renderMessage(msg ChatMessage, styler style.Style, width int) string {
	switch msg.Type {
	case MessageTypeUser:
		return renderUserMessage(msg.Content, styler)
	case MessageTypeAssistant:
		return renderAssistantMessage(msg.Content, styler)
	case MessageTypeScript:
		return renderScript(msg.Content, styler, width)
	case MessageTypeOutput:
		return renderOutput(msg.Content, msg.ExitCode, msg.Duration, msg.ExecutionTime, styler)
	default:
		return msg.Content
	}
}

// renderUserMessage renders a user message with highlighted label.
func renderUserMessage(content string, styler style.Style) string {
	label := lipgloss.NewStyle().
		Bold(true).
		Foreground(styler.HighlightColor().Value()).
		Render("▶ Your Request:")

	text := lipgloss.NewStyle().
		Foreground(styler.TextColor().Value()).
		PaddingLeft(2).
		Render(content)

	return fmt.Sprintf("%s\n%s", label, text)
}

// renderAssistantMessage renders an assistant message with highlighted label.
func renderAssistantMessage(content string, styler style.Style) string {
	label := lipgloss.NewStyle().
		Bold(true).
		Foreground(styler.ActiveColor().Value()).
		Render("[Assistant]:")

	text := lipgloss.NewStyle().
		Foreground(styler.TextColor().Value()).
		Render(content)

	return fmt.Sprintf("%s %s", label, text)
}

// renderOutput renders execution output with muted styling.
func renderOutput(content string, exitCode *int, duration *time.Duration, executionTime *time.Time, styler style.Style) string {
	// Build label parts
	baseLabel := "● Execution Output"

	var labelParts []string
	labelParts = append(labelParts, lipgloss.NewStyle().
		Bold(true).
		Foreground(styler.PlaceholderColor().Value()).
		Render(baseLabel+": "))

	if exitCode != nil && duration != nil && executionTime != nil {
		// Format duration (e.g., "1.2s", "345ms")
		durationStr := formatDuration(*duration)

		// Format execution time (HH:mm:ss)
		timeStr := executionTime.Format("15:04:05")

		// Status with color coding
		var statusStyle lipgloss.Style
		var statusText string
		if *exitCode == 0 {
			statusStyle = lipgloss.NewStyle().
				Foreground(styler.SuccessColor().Value()).
				Bold(true)
			statusText = "✓ Success"
		} else {
			statusStyle = lipgloss.NewStyle().
				Foreground(styler.ErrorColor().Value()).
				Bold(true)
			statusText = fmt.Sprintf("✗ Failed (exit %d)", *exitCode)
		}

		metadataStyle := lipgloss.NewStyle().
			Foreground(styler.PlaceholderColor().Value())

		// Combine: status • duration • time
		labelParts = append(labelParts,
			statusStyle.Render(statusText),
			metadataStyle.Render(" • "),
			metadataStyle.Render(durationStr),
			metadataStyle.Render(" • "),
			metadataStyle.Render(timeStr))
	}

	label := lipgloss.JoinHorizontal(lipgloss.Left, labelParts...)

	// Create styled container with left border accent
	outputStyle := lipgloss.NewStyle().
		Background(styler.VerySubduedColor().Value()).
		Foreground(styler.TextColor().Value()).
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(styler.PlaceholderColor().Value()).
		Padding(1, 2).
		MarginLeft(2).
		MarginTop(1)

	// Strip ANSI escape codes from PTY output before rendering
	cleanContent := stripANSI(strings.TrimRight(content, "\n"))
	styledContent := outputStyle.Render(cleanContent)

	return fmt.Sprintf("%s\n%s", label, styledContent)
}

// stripANSI removes ANSI escape sequences from a string.
func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// renderScript renders a script with syntax highlighting and border.
func renderScript(content string, styler style.Style, width int) string {
	// Add header label
	label := lipgloss.NewStyle().
		Bold(true).
		Foreground(styler.ActiveColor().Value()).
		Render("✓ Generated Script:")

	// Try to syntax highlight the script
	highlighted := highlightCode(content, styler)

	// Create a box style with border and indent
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styler.BorderColor().Value()).
		Padding(0, 1).
		MarginLeft(2)

	// Set max width if provided
	const boxMargin = 4
	if width > boxMargin {
		boxStyle = boxStyle.MaxWidth(width - boxMargin)
	}

	return fmt.Sprintf("%s\n%s", label, boxStyle.Render(highlighted))
}

// highlightCode applies syntax highlighting to code using Chroma.
func highlightCode(code string, styler style.Style) string {
	// Get bash lexer (default for scripts)
	lexer := lexers.Get("bash")
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	// Get terminal formatter
	formatter := formatters.Get("terminal")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	// Get chroma style
	chromaStyle := styles.Get(styler.PreviewColorSchemeName())
	if chromaStyle == nil {
		chromaStyle = styles.Fallback
	}

	// Tokenize and format
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		// If tokenization fails, return plain text
		return code
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, chromaStyle, iterator)
	if err != nil {
		// If formatting fails, return plain text
		return code
	}

	return buf.String()
}

// formatDuration formats a duration into a human-readable string.
func formatDuration(d time.Duration) string {
	const secondsPerMinute = 60
	switch {
	case d < time.Second:
		return fmt.Sprintf("%dms", d.Milliseconds())
	case d < time.Minute:
		return fmt.Sprintf("%.1fs", d.Seconds())
	default:
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % secondsPerMinute
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
}
