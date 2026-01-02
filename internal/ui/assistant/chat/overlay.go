package chat

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// cell represents a single terminal cell with its display content and ANSI styling.
type cell struct {
	content    string // The visible character(s) - may be multi-byte (emoji, CJK chars)
	ansiPrefix string // ANSI codes that apply to this cell (color, bold, etc.)
}

// PlaceOverlay composites a foreground string (e.g., modal) over a background string
// (e.g., chat interface) at the specified position. This preserves the background
// unlike lipgloss.Place which creates a blank canvas.
func PlaceOverlay(
	background string,
	foreground string,
	terminalWidth int,
	terminalHeight int,
	horizontalPos lipgloss.Position,
	verticalPos lipgloss.Position,
) string {
	// Calculate foreground dimensions first
	fgWidth := lipgloss.Width(foreground)
	fgHeight := lipgloss.Height(foreground)

	// Parse both strings into grids (foreground uses its actual size, not terminal size)
	bgGrid := parseToGrid(background, terminalWidth, terminalHeight)
	fgGrid := parseToGrid(foreground, fgWidth, fgHeight)

	// Calculate foreground position
	x, y := calculatePosition(terminalWidth, terminalHeight, fgWidth, fgHeight, horizontalPos, verticalPos)

	// Composite foreground over background
	result := compositeGrids(bgGrid, fgGrid, x, y, fgWidth, fgHeight)

	// Convert back to string
	return gridToString(result)
}

// calculatePosition determines the top-left corner of the overlay based on positioning.
func calculatePosition(
	bgWidth, bgHeight int,
	fgWidth, fgHeight int,
	hPos, vPos lipgloss.Position,
) (x, y int) {
	// Calculate horizontal position
	switch hPos {
	case lipgloss.Left:
		x = 0
	case lipgloss.Center:
		x = (bgWidth - fgWidth) / 2
	case lipgloss.Right:
		x = bgWidth - fgWidth
	}

	// Calculate vertical position
	switch vPos {
	case lipgloss.Top:
		y = 0
	case lipgloss.Center:
		y = (bgHeight - fgHeight) / 2
	case lipgloss.Bottom:
		y = bgHeight - fgHeight
	}

	// Clamp to valid range
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	return x, y
}

// parseToGrid converts a rendered string into a 2D grid of cells.
func parseToGrid(str string, width int, height int) [][]cell {
	lines := strings.Split(str, "\n")
	grid := make([][]cell, height)

	for i := 0; i < height; i++ {
		grid[i] = make([]cell, width)

		if i < len(lines) {
			grid[i] = parseLine(lines[i], width)
		} else {
			// Empty line - fill with empty cells
			for j := 0; j < width; j++ {
				grid[i][j] = cell{content: " ", ansiPrefix: ""}
			}
		}
	}

	return grid
}

// parseLine converts a single line into a row of cells, handling ANSI codes.
func parseLine(line string, width int) []cell {
	row := make([]cell, width)
	currentANSI := ""
	col := 0

	// Process string character by character
	i := 0
	for i < len(line) && col < width {
		// Check if we're at the start of an ANSI sequence
		if line[i] == '\x1b' {
			// Find the complete ANSI sequence using the regex
			remaining := line[i:]
			matches := ansiRegex.FindStringIndex(remaining)
			if matches != nil && matches[0] == 0 {
				// Found ANSI sequence at current position
				ansiSeq := remaining[matches[0]:matches[1]]
				currentANSI += ansiSeq
				i += len(ansiSeq)
				continue
			}
		}

		// Regular character - decode as UTF-8 rune
		r, size := decodeRune(line[i:])
		if r == 0 {
			// Invalid or null character
			i++
			continue
		}

		charWidth := ansi.StringWidth(string(r))

		if charWidth == 0 {
			// Zero-width character, skip
			i += size
			continue
		}

		// Store cell with current ANSI state
		row[col] = cell{
			content:    string(r),
			ansiPrefix: currentANSI,
		}

		col++
		i += size

		// For wide characters (width > 1), fill continuation cells
		for w := 1; w < charWidth && col < width; w++ {
			row[col] = cell{
				content:    "", // Continuation cell
				ansiPrefix: currentANSI,
			}
			col++
		}
	}

	// Fill remaining cells with spaces
	for col < width {
		row[col] = cell{content: " ", ansiPrefix: currentANSI}
		col++
	}

	return row
}

const (
	asciiMax      = 128
	utf8TwoByte   = 0x80
	utf8ThreeByte = 0x800
	utf8FourByte  = 0x10000
	utf8ByteSize1 = 1
	utf8ByteSize2 = 2
	utf8ByteSize3 = 3
	utf8ByteSize4 = 4
)

// decodeRune decodes the first rune from a string, handling UTF-8.
func decodeRune(s string) (rune, int) {
	if len(s) == 0 {
		return 0, 0
	}

	// Check for ASCII fast path
	if s[0] < asciiMax {
		return rune(s[0]), utf8ByteSize1
	}

	// Decode multi-byte UTF-8
	for i, r := range s {
		if i == 0 {
			// Get the byte length of this rune
			size := utf8ByteSize1
			if r >= utf8TwoByte {
				switch {
				case r >= utf8FourByte:
					size = utf8ByteSize4
				case r >= utf8ThreeByte:
					size = utf8ByteSize3
				default:
					size = utf8ByteSize2
				}
			}
			return r, size
		}
	}

	return 0, utf8ByteSize1
}

// compositeGrids overlays foreground grid onto background grid at position (x, y).
func compositeGrids(background [][]cell, foreground [][]cell, x, y, fgWidth, fgHeight int) [][]cell {
	// Create a copy of the background
	result := make([][]cell, len(background))
	for i := range background {
		result[i] = make([]cell, len(background[i]))
		copy(result[i], background[i])
	}

	// Overlay foreground
	for fgY := 0; fgY < fgHeight && fgY < len(foreground); fgY++ {
		bgY := y + fgY
		if bgY < 0 || bgY >= len(result) {
			continue
		}

		for fgX := 0; fgX < fgWidth && fgX < len(foreground[fgY]); fgX++ {
			bgX := x + fgX
			if bgX < 0 || bgX >= len(result[bgY]) {
				continue
			}

			fgCell := foreground[fgY][fgX]
			// Only overlay non-empty cells (this preserves background where modal is transparent)
			// However, modals are typically fully opaque, so we overlay everything
			result[bgY][bgX] = fgCell
		}
	}

	return result
}

// gridToString reconstructs a rendered string from a 2D cell grid.
func gridToString(grid [][]cell) string {
	var builder strings.Builder
	lastANSI := ""

	for rowIdx, row := range grid {
		for _, c := range row {
			// Only emit ANSI codes when they change
			if c.ansiPrefix != lastANSI {
				builder.WriteString(c.ansiPrefix)
				lastANSI = c.ansiPrefix
			}
			if c.content != "" {
				builder.WriteString(c.content)
			}
		}

		// Add newline except for last row
		if rowIdx < len(grid)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}
