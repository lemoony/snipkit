package uimsg

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/pmezard/go-difflib/difflib"

	"github.com/lemoony/snipkit/internal/ui/style"
)

const (
	defaultContextLines        = 3
	minimalConfigLineThreshold = 5
	minTruncateWidth           = 3
	lineNumberWidth            = 4

	// Table column indices.
	colBefore = 0
	colAfter  = 1
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

// createDiffStyleFunc generates a StyleFunc for table cell styling based on diff line types.
func createDiffStyleFunc(diffLines []DiffLine, styler *style.Style) table.StyleFunc {
	return func(row, col int) lipgloss.Style {
		// Header row styling
		if row == -1 {
			return lipgloss.NewStyle().
				Bold(true).
				Foreground(styler.TitleColor().Value()).
				Align(lipgloss.Center)
		}

		// Data rows
		if row >= 0 && row < len(diffLines) {
			dl := diffLines[row]

			// Column 0: BEFORE (left side)
			if col == colBefore {
				return getLeftContentStyle(dl.LineType, styler)
			}

			// Column 1: AFTER (right side)
			if col == colAfter {
				return getRightContentStyle(dl.LineType, styler)
			}
		}

		return lipgloss.NewStyle()
	}
}

// getLeftContentStyle returns the style for left (BEFORE) content based on line type.
func getLeftContentStyle(lineType DiffLineType, styler *style.Style) lipgloss.Style {
	base := lipgloss.NewStyle()
	switch lineType {
	case DiffLineContext:
		return base.Foreground(styler.TextColor().Value())
	case DiffLineRemoved, DiffLineModified:
		return base.Foreground(styler.ErrorColor().Value()).Bold(true)
	case DiffLineAdded:
		return base.Foreground(styler.SubduedColor().Value())
	default:
		return base
	}
}

// getRightContentStyle returns the style for right (AFTER) content based on line type.
func getRightContentStyle(lineType DiffLineType, styler *style.Style) lipgloss.Style {
	base := lipgloss.NewStyle()
	switch lineType {
	case DiffLineContext:
		return base.Foreground(styler.TextColor().Value())
	case DiffLineAdded, DiffLineModified:
		return base.Foreground(styler.SuccessColor().Value()).Bold(true)
	case DiffLineRemoved:
		return base.Foreground(styler.SubduedColor().Value())
	default:
		return base
	}
}

// renderSideBySideTable creates a table-based side-by-side diff view.
func renderSideBySideTable(diffLines []DiffLine, styler *style.Style, width int) string {
	// Create table with 2 columns: BEFORE (line# + content), AFTER (line# + content)
	t := table.New().
		Headers("BEFORE", "AFTER").
		BorderStyle(lipgloss.NewStyle().Foreground(styler.BorderColor().Value())).
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		BorderRow(false).
		Width(width).
		StyleFunc(createDiffStyleFunc(diffLines, styler))

	// Add all diff lines as rows with line numbers prepended to content
	for _, dl := range diffLines {
		leftNum := formatLineNumber(dl.LeftNum, lineNumberWidth)
		rightNum := formatLineNumber(dl.RightNum, lineNumberWidth)
		leftCell := leftNum + dl.LeftLine
		rightCell := rightNum + dl.RightLine
		t.Row(leftCell, rightCell)
	}

	return t.Render()
}

// renderNewConfigOnlyTable renders a single-column table view for new installations.
func renderNewConfigOnlyTable(newYaml string, styler *style.Style, width int) string {
	lines := strings.Split(newYaml, "\n")

	// Create 2-column table: line number and content
	t := table.New().
		Headers("#", "NEW CONFIGURATION").
		BorderStyle(lipgloss.NewStyle().Foreground(styler.SuccessColor().Value())).
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		BorderRow(false).
		Width(width)

	// StyleFunc for new config (all green)
	t.StyleFunc(func(row, col int) lipgloss.Style {
		if row == -1 {
			// Header
			return lipgloss.NewStyle().
				Bold(true).
				Foreground(styler.SuccessColor().Value()).
				Align(lipgloss.Center)
		}

		if col == 0 {
			// Line numbers (right-aligned)
			return lipgloss.NewStyle().
				Foreground(styler.SuccessColor().Value()).
				Align(lipgloss.Right)
		}

		// Content
		return lipgloss.NewStyle().
			Foreground(styler.SuccessColor().Value()).
			Bold(true)
	})

	// Add rows with line numbers
	for i, line := range lines {
		lineNum := fmt.Sprintf("%d", i+1)
		t.Row(lineNum, line)
	}

	return t.Render()
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
