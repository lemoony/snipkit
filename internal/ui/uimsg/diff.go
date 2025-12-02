package uimsg

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pmezard/go-difflib/difflib"

	"github.com/lemoony/snipkit/internal/ui/style"
)

const (
	defaultContextLines        = 3
	minimalConfigLineThreshold = 5
	minTruncateWidth           = 3
)

// DiffLineType represents the type of a diff line.
type DiffLineType int

const (
	DiffLineContext  DiffLineType = iota // Unchanged line (both sides)
	DiffLineAdded                        // Added line (right only, green)
	DiffLineRemoved                      // Removed line (left only, red)
	DiffLineModified                     // Modified line (both sides, different)
)

// DiffLine represents a single line in the diff view.
type DiffLine struct {
	LineType  DiffLineType
	LeftLine  string // Line content for left side (old)
	RightLine string // Line content for right side (new)
	LeftNum   int    // Line number on left (0 if not applicable)
	RightNum  int    // Line number on right (0 if not applicable)
}

// computeDiff compares two YAML strings and returns a slice of DiffLines.
func computeDiff(oldYaml, newYaml string) []DiffLine {
	oldLines := strings.Split(oldYaml, "\n")
	newLines := strings.Split(newYaml, "\n")

	matcher := difflib.NewMatcher(oldLines, newLines)
	opcodes := matcher.GetOpCodes()

	diffLines := processOpcodes(opcodes, oldLines, newLines)

	// Apply context compression
	return compressContext(diffLines, defaultContextLines)
}

// processOpcodes converts difflib opcodes into DiffLines.
func processOpcodes(opcodes []difflib.OpCode, oldLines, newLines []string) []DiffLine {
	var diffLines []DiffLine
	oldLineNum, newLineNum := 1, 1

	for _, opcode := range opcodes {
		var lines []DiffLine
		switch opcode.Tag {
		case 'e': // equal - show context
			lines, oldLineNum, newLineNum = processEqualLines(opcode, oldLines, newLines, oldLineNum, newLineNum)
		case 'r': // replace - show both sides as modified
			lines, oldLineNum, newLineNum = processReplaceLines(opcode, oldLines, newLines, oldLineNum, newLineNum)
		case 'd': // delete - show on left only
			lines, oldLineNum = processDeleteLines(opcode, oldLines, oldLineNum)
		case 'i': // insert - show on right only
			lines, newLineNum = processInsertLines(opcode, newLines, newLineNum)
		}
		diffLines = append(diffLines, lines...)
	}

	return diffLines
}

func processEqualLines(opcode difflib.OpCode, oldLines, newLines []string, oldLineNum, newLineNum int) ([]DiffLine, int, int) {
	var result []DiffLine
	i1, i2, j1 := opcode.I1, opcode.I2, opcode.J1

	for i := i1; i < i2; i++ {
		result = append(result, DiffLine{
			LineType:  DiffLineContext,
			LeftLine:  oldLines[i],
			RightLine: newLines[j1+(i-i1)],
			LeftNum:   oldLineNum,
			RightNum:  newLineNum,
		})
		oldLineNum++
		newLineNum++
	}
	return result, oldLineNum, newLineNum
}

func processReplaceLines(opcode difflib.OpCode, oldLines, newLines []string, oldLineNum, newLineNum int) ([]DiffLine, int, int) {
	var result []DiffLine
	i1, i2, j1, j2 := opcode.I1, opcode.I2, opcode.J1, opcode.J2

	maxLen := i2 - i1
	if j2-j1 > maxLen {
		maxLen = j2 - j1
	}

	for k := 0; k < maxLen; k++ {
		leftLine, rightLine := "", ""
		leftNum, rightNum := 0, 0

		if i1+k < i2 {
			leftLine = oldLines[i1+k]
			leftNum = oldLineNum
			oldLineNum++
		}
		if j1+k < j2 {
			rightLine = newLines[j1+k]
			rightNum = newLineNum
			newLineNum++
		}

		result = append(result, DiffLine{
			LineType:  DiffLineModified,
			LeftLine:  leftLine,
			RightLine: rightLine,
			LeftNum:   leftNum,
			RightNum:  rightNum,
		})
	}
	return result, oldLineNum, newLineNum
}

func processDeleteLines(opcode difflib.OpCode, oldLines []string, oldLineNum int) ([]DiffLine, int) {
	var result []DiffLine
	for i := opcode.I1; i < opcode.I2; i++ {
		result = append(result, DiffLine{
			LineType:  DiffLineRemoved,
			LeftLine:  oldLines[i],
			RightLine: "",
			LeftNum:   oldLineNum,
			RightNum:  0,
		})
		oldLineNum++
	}
	return result, oldLineNum
}

func processInsertLines(opcode difflib.OpCode, newLines []string, newLineNum int) ([]DiffLine, int) {
	var result []DiffLine
	for j := opcode.J1; j < opcode.J2; j++ {
		result = append(result, DiffLine{
			LineType:  DiffLineAdded,
			LeftLine:  "",
			RightLine: newLines[j],
			LeftNum:   0,
			RightNum:  newLineNum,
		})
		newLineNum++
	}
	return result, newLineNum
}

// compressContext reduces long stretches of unchanged lines to show only
// N lines of context before and after each change.
func compressContext(diffLines []DiffLine, contextLines int) []DiffLine {
	if len(diffLines) == 0 {
		return diffLines
	}

	// Find indices of all changed lines
	changeIndices := []int{}
	for i, line := range diffLines {
		if line.LineType != DiffLineContext {
			changeIndices = append(changeIndices, i)
		}
	}

	// If no changes, return as-is (shouldn't happen in practice)
	if len(changeIndices) == 0 {
		return diffLines
	}

	// Calculate which context lines to keep
	keepLines := make(map[int]bool)
	for _, idx := range changeIndices {
		// Keep the changed line itself
		keepLines[idx] = true
		// Keep context before
		for i := idx - contextLines; i < idx; i++ {
			if i >= 0 {
				keepLines[i] = true
			}
		}
		// Keep context after
		for i := idx + 1; i <= idx+contextLines; i++ {
			if i < len(diffLines) {
				keepLines[i] = true
			}
		}
	}

	// Build result with compression markers
	var result []DiffLine
	skipping := false
	skippedCount := 0

	for i, line := range diffLines {
		if keepLines[i] {
			// If we were skipping, add a separator
			if skipping {
				result = append(result, DiffLine{
					LineType:  DiffLineContext,
					LeftLine:  fmt.Sprintf("... (%d lines unchanged)", skippedCount),
					RightLine: fmt.Sprintf("... (%d lines unchanged)", skippedCount),
					LeftNum:   0,
					RightNum:  0,
				})
				skipping = false
				skippedCount = 0
			}
			result = append(result, line)
		} else {
			skipping = true
			skippedCount++
		}
	}

	return result
}

// renderHeader creates the header and separator line for the diff view.
func renderHeader(leftWidth, rightWidth int, styler *style.Style, borderStyle lipgloss.Style) []string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styler.TitleColor().Value())

	// Build header manually with proper padding
	leftHeaderText := "BEFORE"
	rightHeaderText := "AFTER"
	leftPadding := (leftWidth - len(leftHeaderText)) / 2
	rightPadding := (rightWidth - len(rightHeaderText)) / 2

	leftHeader := headerStyle.Render(strings.Repeat(" ", leftPadding) + leftHeaderText + strings.Repeat(" ", leftWidth-leftPadding-len(leftHeaderText)))
	midSeparator := borderStyle.Render(" │ ")
	rightHeader := headerStyle.Render(strings.Repeat(" ", rightPadding) + rightHeaderText + strings.Repeat(" ", rightWidth-rightPadding-len(rightHeaderText)))

	header := leftHeader + midSeparator + rightHeader

	// Add separator line
	separatorLine := strings.Repeat("─", leftWidth) + "─┼─" + strings.Repeat("─", rightWidth)

	return []string{header, borderStyle.Render(separatorLine)}
}

// renderSideBySide creates a styled side-by-side diff view.
func renderSideBySide(diffLines []DiffLine, styler *style.Style, width int) string {
	// Width allocation
	lineNumWidth := 4
	separatorWidth := 3
	contentWidth := (width - (lineNumWidth * 2) - separatorWidth) / 2

	leftWidth := lineNumWidth + contentWidth
	rightWidth := lineNumWidth + contentWidth

	// Border style for separators
	borderStyle := lipgloss.NewStyle().Foreground(styler.BorderColor().Value())

	var lines []string
	lines = append(lines, renderHeader(leftWidth, rightWidth, styler, borderStyle)...)

	// Render each diff line
	for _, dl := range diffLines {
		leftStyle, rightStyle := getStylesForLineType(dl.LineType, styler)

		// Format line numbers
		leftNumStr := formatLineNumber(dl.LeftNum, lineNumWidth)
		rightNumStr := formatLineNumber(dl.RightNum, lineNumWidth)

		// Truncate content if needed
		leftContent := truncateString(dl.LeftLine, contentWidth)
		rightContent := truncateString(dl.RightLine, contentWidth)

		// Apply styles and combine
		leftPart := lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftStyle.Width(lineNumWidth).Render(leftNumStr),
			leftStyle.Width(contentWidth).Render(leftContent),
		)

		sep := lipgloss.NewStyle().
			Foreground(styler.BorderColor().Value()).
			Width(separatorWidth).
			Render(" │ ")

		rightPart := lipgloss.JoinHorizontal(
			lipgloss.Top,
			rightStyle.Width(lineNumWidth).Render(rightNumStr),
			rightStyle.Width(contentWidth).Render(rightContent),
		)

		line := lipgloss.JoinHorizontal(lipgloss.Top, leftPart, sep, rightPart)
		lines = append(lines, line)
	}

	return wrapWithBorder(lines, borderStyle, width)
}

// wrapWithBorder adds a full box border around the diff content.
func wrapWithBorder(lines []string, borderStyle lipgloss.Style, width int) string {
	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Add left and right borders with padding
	var borderedLines []string
	for _, line := range strings.Split(content, "\n") {
		borderedLine := borderStyle.Render("│") + " " + line + " " + borderStyle.Render("│")
		borderedLines = append(borderedLines, borderedLine)
	}

	// Add top and bottom borders
	topBorder := borderStyle.Render("┌" + strings.Repeat("─", width+2) + "┐")
	bottomBorder := borderStyle.Render("└" + strings.Repeat("─", width+2) + "┘")

	return topBorder + "\n" + strings.Join(borderedLines, "\n") + "\n" + bottomBorder
}

// getStylesForLineType returns appropriate lipgloss styles for left and right sides.
func getStylesForLineType(lineType DiffLineType, styler *style.Style) (left, right lipgloss.Style) {
	base := lipgloss.NewStyle()

	switch lineType {
	case DiffLineContext:
		return base.Foreground(styler.TextColor().Value()),
			base.Foreground(styler.TextColor().Value())

	case DiffLineAdded:
		return base.Foreground(styler.SubduedColor().Value()),
			base.Foreground(styler.SuccessColor().Value()).Bold(true)

	case DiffLineRemoved:
		return base.Foreground(styler.ErrorColor().Value()).Bold(true),
			base.Foreground(styler.SubduedColor().Value())

	case DiffLineModified:
		return base.Foreground(styler.ErrorColor().Value()),
			base.Foreground(styler.SuccessColor().Value())

	default:
		return base, base
	}
}

// formatLineNumber formats a line number for display.
func formatLineNumber(num int, width int) string {
	if num == 0 {
		return strings.Repeat(" ", width)
	}
	return fmt.Sprintf("%*d ", width-1, num)
}

// truncateString truncates a string to fit within maxWidth.
func truncateString(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth <= minTruncateWidth {
		return s[:maxWidth]
	}
	return s[:maxWidth-minTruncateWidth] + "..."
}

// isMinimalConfig detects if a config is essentially empty or minimal.
func isMinimalConfig(configStr string) bool {
	lines := strings.Split(configStr, "\n")
	meaningfulLines := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			meaningfulLines++
		}
	}

	// If less than threshold meaningful lines, consider it minimal
	return meaningfulLines < minimalConfigLineThreshold
}

// renderNewConfigOnly renders a single-column view for new installations.
func renderNewConfigOnly(newYaml string, styler *style.Style, width int) string {
	lines := strings.Split(newYaml, "\n")
	var styledLines []string

	addedStyle := lipgloss.NewStyle().
		Foreground(styler.SuccessColor().Value()).
		Bold(true)

	for i, line := range lines {
		lineNum := fmt.Sprintf("%4d + ", i+1)
		styledLine := addedStyle.Render(lineNum + line)
		styledLines = append(styledLines, styledLine)
	}

	blockStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styler.SuccessColor().Value()).
		Padding(1).
		Width(width)

	content := lipgloss.JoinVertical(lipgloss.Left, styledLines...)
	return blockStyle.Render(styler.Title("NEW CONFIGURATION") + "\n\n" + content)
}
