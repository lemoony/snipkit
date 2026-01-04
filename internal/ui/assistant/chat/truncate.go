package chat

import (
	"fmt"
	"strings"
)

const maxLines = 20

// truncateContent truncates content to maxLines, adding an indicator if truncated.
func truncateContent(content string, maxLines int) string {
	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines {
		return content
	}

	truncatedLines := lines[:maxLines]
	remaining := len(lines) - maxLines
	result := strings.Join(truncatedLines, "\n")
	result += fmt.Sprintf("\n\n... (%d more lines)", remaining)
	return result
}
