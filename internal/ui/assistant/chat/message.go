package chat

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/lipgloss"

	"github.com/lemoony/snipkit/internal/ui/style"
)

type MessageType int

const (
	MessageTypeUser MessageType = iota
	MessageTypeAssistant
	MessageTypeScript
	MessageTypeOutput
)

type ChatMessage struct {
	Type      MessageType
	Content   string
	Timestamp time.Time
}

const maxMessagesPerEntry = 3 // User prompt, script, and output

// buildMessagesFromHistory converts a history of interactions into chat messages.
func buildMessagesFromHistory(history []HistoryEntry) []ChatMessage {
	if len(history) == 0 {
		return []ChatMessage{}
	}

	messages := make([]ChatMessage, 0, len(history)*maxMessagesPerEntry)
	for _, entry := range history {
		// Add user prompt
		messages = append(messages, ChatMessage{
			Type:      MessageTypeUser,
			Content:   entry.UserPrompt,
			Timestamp: time.Now(),
		})

		// Add generated script (truncated)
		if entry.GeneratedScript != "" {
			script := truncateContent(entry.GeneratedScript, maxLines)
			messages = append(messages, ChatMessage{
				Type:      MessageTypeScript,
				Content:   script,
				Timestamp: time.Now(),
			})
		}

		// Add execution output (truncated)
		if entry.ExecutionOutput != "" {
			output := truncateContent(entry.ExecutionOutput, maxLines)
			messages = append(messages, ChatMessage{
				Type:      MessageTypeOutput,
				Content:   output,
				Timestamp: time.Now(),
			})
		}
	}

	return messages
}

// renderMessages renders all messages into a single string for the viewport.
func renderMessages(messages []ChatMessage, styler style.Style, width int) string {
	if len(messages) == 0 {
		contextNote := styler.PromptDescription(
			"Your prompts and their results are automatically provided as context to the AI.",
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
		return renderOutput(msg.Content, styler)
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
func renderOutput(content string, styler style.Style) string {
	label := lipgloss.NewStyle().
		Bold(true).
		Foreground(styler.PlaceholderColor().Value()).
		Render("● Execution Output:")

	// Create styled container with left border accent
	outputStyle := lipgloss.NewStyle().
		Background(styler.VerySubduedColor().Value()).
		Foreground(styler.PlaceholderColor().Value()).
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(styler.PlaceholderColor().Value()).
		Padding(1, 2).
		MarginLeft(2).
		MarginTop(1)

	styledContent := outputStyle.Render(content)

	return fmt.Sprintf("%s\n%s", label, styledContent)
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
