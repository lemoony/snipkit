package chat

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

const (
	testBackground4x4 = "AAAA\nAAAA\nAAAA\nAAAA"
	testForeground2x2 = "BB\nBB"
)

func Test_PlaceOverlay_CenterPosition(t *testing.T) {
	background := testBackground4x4
	foreground := testForeground2x2

	result := PlaceOverlay(background, foreground, 4, 4, lipgloss.Center, lipgloss.Center)

	lines := strings.Split(result, "\n")
	assert.Len(t, lines, 4)
	// Centered: row 1 and 2 should contain BB in the middle
	assert.Contains(t, lines[1], "BB")
	assert.Contains(t, lines[2], "BB")
}

func Test_PlaceOverlay_TopLeft(t *testing.T) {
	background := testBackground4x4
	foreground := testForeground2x2

	result := PlaceOverlay(background, foreground, 4, 4, lipgloss.Left, lipgloss.Top)

	lines := strings.Split(result, "\n")
	assert.Len(t, lines, 4)
	// Top-left: rows 0 and 1 should start with BB
	assert.True(t, strings.HasPrefix(lines[0], "BB"))
	assert.True(t, strings.HasPrefix(lines[1], "BB"))
}

func Test_PlaceOverlay_BottomRight(t *testing.T) {
	background := testBackground4x4
	foreground := testForeground2x2

	result := PlaceOverlay(background, foreground, 4, 4, lipgloss.Right, lipgloss.Bottom)

	lines := strings.Split(result, "\n")
	assert.Len(t, lines, 4)
	// Bottom-right: rows 2 and 3 should end with BB
	assert.True(t, strings.HasSuffix(lines[2], "BB"))
	assert.True(t, strings.HasSuffix(lines[3], "BB"))
}

func Test_PlaceOverlay_EmptyForeground(t *testing.T) {
	background := "AAAA\nAAAA"
	foreground := ""

	result := PlaceOverlay(background, foreground, 4, 2, lipgloss.Center, lipgloss.Center)

	// Should return background unchanged
	assert.NotEmpty(t, result)
}

func Test_PlaceOverlay_EmptyBackground(t *testing.T) {
	background := ""
	foreground := "BB"

	// Empty background with terminal dimensions
	result := PlaceOverlay(background, foreground, 4, 2, lipgloss.Center, lipgloss.Center)

	// Should still produce output
	assert.NotEmpty(t, result)
}

func Test_PlaceOverlay_OversizedForeground(t *testing.T) {
	background := "AA\nAA"
	foreground := "BBBBBB\nBBBBBB\nBBBBBB\nBBBBBB"

	// Foreground larger than background - should clamp to background size
	result := PlaceOverlay(background, foreground, 2, 2, lipgloss.Center, lipgloss.Center)

	assert.NotEmpty(t, result)
	lines := strings.Split(result, "\n")
	assert.Equal(t, 2, len(lines))
}

func Test_PlaceOverlay_AllPositions(t *testing.T) {
	tests := []struct {
		name   string
		hPos   lipgloss.Position
		vPos   lipgloss.Position
		checkX int // expected column for foreground
		checkY int // expected row for foreground
	}{
		{"top-left", lipgloss.Left, lipgloss.Top, 0, 0},
		{"top-center", lipgloss.Center, lipgloss.Top, 1, 0},
		{"top-right", lipgloss.Right, lipgloss.Top, 2, 0},
		{"center-left", lipgloss.Left, lipgloss.Center, 0, 1},
		{"center-center", lipgloss.Center, lipgloss.Center, 1, 1},
		{"center-right", lipgloss.Right, lipgloss.Center, 2, 1},
		{"bottom-left", lipgloss.Left, lipgloss.Bottom, 0, 2},
		{"bottom-center", lipgloss.Center, lipgloss.Bottom, 1, 2},
		{"bottom-right", lipgloss.Right, lipgloss.Bottom, 2, 2},
	}

	background := "AAA\nAAA\nAAA"
	foreground := "B"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PlaceOverlay(background, foreground, 3, 3, tt.hPos, tt.vPos)
			lines := strings.Split(result, "\n")

			assert.Len(t, lines, 3)
			// Check that B is at the expected position
			assert.Equal(t, "B", string(lines[tt.checkY][tt.checkX]))
		})
	}
}

func Test_calculatePosition(t *testing.T) {
	tests := []struct {
		name      string
		bgWidth   int
		bgHeight  int
		fgWidth   int
		fgHeight  int
		hPos      lipgloss.Position
		vPos      lipgloss.Position
		expectedX int
		expectedY int
	}{
		{"left-top", 10, 10, 4, 4, lipgloss.Left, lipgloss.Top, 0, 0},
		{"center-center", 10, 10, 4, 4, lipgloss.Center, lipgloss.Center, 3, 3},
		{"right-bottom", 10, 10, 4, 4, lipgloss.Right, lipgloss.Bottom, 6, 6},
		{"oversized-width clamps", 4, 10, 10, 4, lipgloss.Center, lipgloss.Center, 0, 3},
		{"oversized-height clamps", 10, 4, 4, 10, lipgloss.Center, lipgloss.Center, 3, 0},
		{"equal sizes", 5, 5, 5, 5, lipgloss.Center, lipgloss.Center, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := calculatePosition(tt.bgWidth, tt.bgHeight, tt.fgWidth, tt.fgHeight, tt.hPos, tt.vPos)
			assert.Equal(t, tt.expectedX, x)
			assert.Equal(t, tt.expectedY, y)
		})
	}
}

func Test_parseToGrid_BasicString(t *testing.T) {
	str := "ABC\nDEF"
	grid := parseToGrid(str, 3, 2)

	assert.Len(t, grid, 2)
	assert.Len(t, grid[0], 3)

	assert.Equal(t, "A", grid[0][0].content)
	assert.Equal(t, "B", grid[0][1].content)
	assert.Equal(t, "C", grid[0][2].content)
	assert.Equal(t, "D", grid[1][0].content)
	assert.Equal(t, "E", grid[1][1].content)
	assert.Equal(t, "F", grid[1][2].content)
}

func Test_parseToGrid_PadsShortLines(t *testing.T) {
	str := "A\nBC"
	grid := parseToGrid(str, 3, 2)

	assert.Len(t, grid, 2)
	assert.Len(t, grid[0], 3)

	assert.Equal(t, "A", grid[0][0].content)
	assert.Equal(t, " ", grid[0][1].content)
	assert.Equal(t, " ", grid[0][2].content)
}

func Test_parseToGrid_PadsExtraRows(t *testing.T) {
	str := "AB"
	grid := parseToGrid(str, 2, 3)

	assert.Len(t, grid, 3)
	// First row should have content
	assert.Equal(t, "A", grid[0][0].content)
	assert.Equal(t, "B", grid[0][1].content)
	// Extra rows should be spaces
	assert.Equal(t, " ", grid[1][0].content)
	assert.Equal(t, " ", grid[2][0].content)
}

func Test_parseLine_ASCII(t *testing.T) {
	line := "Hello"
	row := parseLine(line, 5)

	assert.Len(t, row, 5)
	assert.Equal(t, "H", row[0].content)
	assert.Equal(t, "e", row[1].content)
	assert.Equal(t, "l", row[2].content)
	assert.Equal(t, "l", row[3].content)
	assert.Equal(t, "o", row[4].content)
}

func Test_parseLine_PadsToWidth(t *testing.T) {
	line := "Hi"
	row := parseLine(line, 5)

	assert.Len(t, row, 5)
	assert.Equal(t, "H", row[0].content)
	assert.Equal(t, "i", row[1].content)
	assert.Equal(t, " ", row[2].content)
	assert.Equal(t, " ", row[3].content)
	assert.Equal(t, " ", row[4].content)
}

func Test_parseLine_WithANSI(t *testing.T) {
	line := "\x1b[31mRED\x1b[0m"
	row := parseLine(line, 5)

	assert.Len(t, row, 5)
	assert.Equal(t, "R", row[0].content)
	assert.Contains(t, row[0].ansiPrefix, "\x1b[31m")
	assert.Equal(t, "E", row[1].content)
	assert.Equal(t, "D", row[2].content)
}

func Test_parseLine_TruncatesLongLine(t *testing.T) {
	line := "Hello World"
	row := parseLine(line, 5)

	assert.Len(t, row, 5)
	assert.Equal(t, "H", row[0].content)
	assert.Equal(t, "e", row[1].content)
	assert.Equal(t, "l", row[2].content)
	assert.Equal(t, "l", row[3].content)
	assert.Equal(t, "o", row[4].content)
}

func Test_parseLine_UTF8MultibyteCharacters(t *testing.T) {
	// Test with UTF-8 characters (Japanese)
	line := "ABこ"
	row := parseLine(line, 5)

	assert.Len(t, row, 5)
	assert.Equal(t, "A", row[0].content)
	assert.Equal(t, "B", row[1].content)
	// Japanese character takes 2 columns, so position 2 has the char and 3 is continuation
	assert.Equal(t, "こ", row[2].content)
	assert.Equal(t, "", row[3].content) // Continuation cell for wide char
}

func Test_decodeRune_ASCII(t *testing.T) {
	tests := []struct {
		input    string
		expected rune
		size     int
	}{
		{"A", 'A', 1},
		{"z", 'z', 1},
		{"1", '1', 1},
		{" ", ' ', 1},
	}

	for _, tt := range tests {
		t.Run(string(tt.expected), func(t *testing.T) {
			r, size := decodeRune(tt.input)
			assert.Equal(t, tt.expected, r)
			assert.Equal(t, tt.size, size)
		})
	}
}

func Test_decodeRune_EmptyString(t *testing.T) {
	r, size := decodeRune("")
	assert.Equal(t, rune(0), r)
	assert.Equal(t, 0, size)
}

func Test_decodeRune_Multibyte(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected rune
	}{
		{"2-byte UTF-8", "é", 'é'},
		{"3-byte UTF-8", "日", '日'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, size := decodeRune(tt.input)
			assert.Equal(t, tt.expected, r)
			assert.Greater(t, size, 1)
		})
	}
}

func Test_compositeGrids_Basic(t *testing.T) {
	// Create 3x3 background of A's
	bg := [][]cell{
		{{content: "A"}, {content: "A"}, {content: "A"}},
		{{content: "A"}, {content: "A"}, {content: "A"}},
		{{content: "A"}, {content: "A"}, {content: "A"}},
	}

	// Create 1x1 foreground of B
	fg := [][]cell{
		{{content: "B"}},
	}

	// Place at center (1, 1)
	result := compositeGrids(bg, fg, 1, 1, 1, 1)

	// Check result
	assert.Equal(t, "A", result[0][0].content)
	assert.Equal(t, "A", result[0][1].content)
	assert.Equal(t, "A", result[1][0].content)
	assert.Equal(t, "B", result[1][1].content) // Overlaid cell
	assert.Equal(t, "A", result[1][2].content)
	assert.Equal(t, "A", result[2][2].content)
}

func Test_compositeGrids_DoesNotModifyOriginal(t *testing.T) {
	bg := [][]cell{
		{{content: "A"}, {content: "A"}},
		{{content: "A"}, {content: "A"}},
	}
	fg := [][]cell{
		{{content: "B"}},
	}

	_ = compositeGrids(bg, fg, 0, 0, 1, 1)

	// Original should be unchanged
	assert.Equal(t, "A", bg[0][0].content)
}

func Test_compositeGrids_OutOfBounds(t *testing.T) {
	bg := [][]cell{
		{{content: "A"}, {content: "A"}},
	}
	fg := [][]cell{
		{{content: "B"}, {content: "B"}, {content: "B"}},
	}

	// Foreground extends beyond background - should not panic
	result := compositeGrids(bg, fg, 0, 0, 3, 1)

	// Should only overlay what fits
	assert.Equal(t, "B", result[0][0].content)
	assert.Equal(t, "B", result[0][1].content)
}

func Test_gridToString_Basic(t *testing.T) {
	grid := [][]cell{
		{{content: "A"}, {content: "B"}},
		{{content: "C"}, {content: "D"}},
	}

	result := gridToString(grid)

	assert.Equal(t, "AB\nCD", result)
}

func Test_gridToString_WithANSI(t *testing.T) {
	grid := [][]cell{
		{{content: "A", ansiPrefix: "\x1b[31m"}, {content: "B", ansiPrefix: "\x1b[31m"}},
	}

	result := gridToString(grid)

	// Should contain ANSI code and both chars
	assert.Contains(t, result, "\x1b[31m")
	assert.Contains(t, result, "A")
	assert.Contains(t, result, "B")
}

func Test_gridToString_SkipsContinuationCells(t *testing.T) {
	grid := [][]cell{
		{{content: "日"}, {content: ""}}, // Wide char with continuation
	}

	result := gridToString(grid)

	assert.Equal(t, "日", result)
}

func Test_gridToString_NoTrailingNewline(t *testing.T) {
	grid := [][]cell{
		{{content: "A"}},
		{{content: "B"}},
	}

	result := gridToString(grid)

	assert.Equal(t, "A\nB", result)
	assert.False(t, strings.HasSuffix(result, "\n"))
}

func Test_gridToString_OptimizesANSICodes(t *testing.T) {
	// Same ANSI code repeated - should only appear once
	ansi := "\x1b[32m"
	grid := [][]cell{
		{{content: "A", ansiPrefix: ansi}, {content: "B", ansiPrefix: ansi}, {content: "C", ansiPrefix: ansi}},
	}

	result := gridToString(grid)

	// Count occurrences - ANSI code should appear only once at the start
	count := strings.Count(result, ansi)
	assert.Equal(t, 1, count)
}
