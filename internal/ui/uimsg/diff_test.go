package uimsg

import (
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
)

const testConfigV110 = "version: v1.1.0\nconfig:\n  editor: vim"

func TestComputeDiff_IdenticalConfigs(t *testing.T) {
	old := testConfigV110
	new := testConfigV110

	diff := computeDiff(old, new)

	// All lines should be context
	for _, line := range diff {
		assert.Equal(t, DiffLineContext, line.LineType)
	}
}

func TestComputeDiff_OnlyAdditions(t *testing.T) {
	old := "version: v1.1.0"
	new := "version: v1.1.0\nconfig:\n  editor: vim"

	diff := computeDiff(old, new)

	hasAdded := false
	for _, line := range diff {
		if line.LineType == DiffLineAdded {
			hasAdded = true
			assert.Empty(t, line.LeftLine)
			assert.NotEmpty(t, line.RightLine)
			assert.Equal(t, 0, line.LeftNum)
			assert.NotEqual(t, 0, line.RightNum)
		}
	}
	assert.True(t, hasAdded, "Should have at least one added line")
}

func TestComputeDiff_OnlyDeletions(t *testing.T) {
	old := testConfigV110
	new := "version: v1.1.0"

	diff := computeDiff(old, new)

	hasRemoved := false
	for _, line := range diff {
		if line.LineType == DiffLineRemoved {
			hasRemoved = true
			assert.NotEmpty(t, line.LeftLine)
			assert.Empty(t, line.RightLine)
			assert.NotEqual(t, 0, line.LeftNum)
			assert.Equal(t, 0, line.RightNum)
		}
	}
	assert.True(t, hasRemoved, "Should have at least one removed line")
}

func TestComputeDiff_Modifications(t *testing.T) {
	old := testConfigV110
	new := "version: v1.3.0\nconfig:\n  editor: nvim"

	diff := computeDiff(old, new)

	hasModified := false
	for _, line := range diff {
		if line.LineType == DiffLineModified || line.LineType == DiffLineRemoved || line.LineType == DiffLineAdded {
			hasModified = true
		}
	}
	assert.True(t, hasModified, "Should have at least one modified/added/removed line")
}

func TestIsMinimalConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   string
		expected bool
	}{
		{
			name:     "empty config",
			config:   "",
			expected: true,
		},
		{
			name:     "minimal config",
			config:   "version: v1.1.0\nconfig: {}",
			expected: true,
		},
		{
			name:     "config with comments only",
			config:   "# This is a comment\n# Another comment\nversion: v1.1.0",
			expected: true,
		},
		{
			name: "full config",
			config: `version: v1.1.0
config:
  editor: vim
  fuzzySearch: true
  manager:
    fslibrary:
      path: ~/snippets`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isMinimalConfig(tt.config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatLineNumber(t *testing.T) {
	tests := []struct {
		num      int
		width    int
		expected string
	}{
		{0, 4, "    "},
		{1, 4, "  1 "},
		{10, 4, " 10 "},
		{100, 4, "100 "},
		{1, 5, "   1 "},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := formatLineNumber(tt.num, tt.width)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.width, len(result))
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxWidth int
		expected string
	}{
		{
			name:     "no truncation needed",
			input:    "short",
			maxWidth: 10,
			expected: "short",
		},
		{
			name:     "exact length",
			input:    "exactly10c",
			maxWidth: 10,
			expected: "exactly10c",
		},
		{
			name:     "truncation needed",
			input:    "this is a very long string",
			maxWidth: 10,
			expected: "this is...",
		},
		{
			name:     "truncation with small width",
			input:    "toolong",
			maxWidth: 3,
			expected: "too",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxWidth)
			assert.Equal(t, tt.expected, result)
			assert.LessOrEqual(t, len(result), tt.maxWidth)
		})
	}
}

func TestRenderSideBySideTable(t *testing.T) {
	// Test that rendering doesn't panic and produces output
	diffLines := []DiffLine{
		{LineType: DiffLineContext, LeftLine: "version: v1.1.0", RightLine: "version: v1.1.0", LeftNum: 1, RightNum: 1},
		{LineType: DiffLineRemoved, LeftLine: "config:", RightLine: "", LeftNum: 2, RightNum: 0},
		{LineType: DiffLineAdded, LeftLine: "", RightLine: "  editor: nvim", LeftNum: 0, RightNum: 2},
	}

	styler := testStyle // Use existing test style from uimsg_test.go
	result := renderSideBySideTable(diffLines, styler, 120)

	assert.NotEmpty(t, result)
	// Verify it contains expected elements (after stripping ANSI codes would be better, but basic check)
	assert.Contains(t, result, "BEFORE")
	assert.Contains(t, result, "AFTER")
}

func TestRenderNewConfigOnlyTable(t *testing.T) {
	newYaml := testConfigV110
	styler := testStyle

	result := renderNewConfigOnlyTable(newYaml, styler, 120)

	assert.NotEmpty(t, result)
	assert.Contains(t, result, "NEW CONFIGURATION")
}

func TestCompressContext(t *testing.T) {
	// Create a diff with a long stretch of context
	diffLines := []DiffLine{
		{LineType: DiffLineContext, LeftLine: "line1", RightLine: "line1", LeftNum: 1, RightNum: 1},
		{LineType: DiffLineContext, LeftLine: "line2", RightLine: "line2", LeftNum: 2, RightNum: 2},
		{LineType: DiffLineContext, LeftLine: "line3", RightLine: "line3", LeftNum: 3, RightNum: 3},
		{LineType: DiffLineContext, LeftLine: "line4", RightLine: "line4", LeftNum: 4, RightNum: 4},
		{LineType: DiffLineContext, LeftLine: "line5", RightLine: "line5", LeftNum: 5, RightNum: 5},
		{LineType: DiffLineContext, LeftLine: "line6", RightLine: "line6", LeftNum: 6, RightNum: 6},
		{LineType: DiffLineContext, LeftLine: "line7", RightLine: "line7", LeftNum: 7, RightNum: 7},
		{LineType: DiffLineContext, LeftLine: "line8", RightLine: "line8", LeftNum: 8, RightNum: 8},
		{LineType: DiffLineAdded, LeftLine: "", RightLine: "new line", LeftNum: 0, RightNum: 9},
		{LineType: DiffLineContext, LeftLine: "line9", RightLine: "line10", LeftNum: 9, RightNum: 10},
		{LineType: DiffLineContext, LeftLine: "line10", RightLine: "line11", LeftNum: 10, RightNum: 11},
	}

	result := compressContext(diffLines, 2)

	// Should be shorter than original due to compression
	assert.Less(t, len(result), len(diffLines))

	// Should contain a compression marker
	hasCompressionMarker := false
	for _, line := range result {
		if strings.Contains(line.LeftLine, "lines unchanged") {
			hasCompressionMarker = true
			break
		}
	}
	assert.True(t, hasCompressionMarker, "Should have compression marker")
}

func Test_getLeftContentStyle_AllLineTypes(t *testing.T) {
	styler := testStyle

	tests := []struct {
		name     string
		lineType DiffLineType
	}{
		{name: "context line", lineType: DiffLineContext},
		{name: "removed line", lineType: DiffLineRemoved},
		{name: "modified line", lineType: DiffLineModified},
		{name: "added line", lineType: DiffLineAdded},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := getLeftContentStyle(tt.lineType, styler)
			assert.NotNil(t, style)
			// Style should be returned without panicking
		})
	}
}

func Test_getRightContentStyle_AllLineTypes(t *testing.T) {
	styler := testStyle

	tests := []struct {
		name     string
		lineType DiffLineType
	}{
		{name: "context line", lineType: DiffLineContext},
		{name: "added line", lineType: DiffLineAdded},
		{name: "modified line", lineType: DiffLineModified},
		{name: "removed line", lineType: DiffLineRemoved},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := getRightContentStyle(tt.lineType, styler)
			assert.NotNil(t, style)
			// Style should be returned without panicking
		})
	}
}

func Test_processReplaceLines_UnequalSides(t *testing.T) {
	tests := []struct {
		name     string
		oldLines []string
		newLines []string
		opcode   difflib.OpCode
		verify   func(t *testing.T, result []DiffLine)
	}{
		{
			name:     "more old lines than new",
			oldLines: []string{"line1", "line2", "line3"},
			newLines: []string{"new1"},
			opcode:   difflib.OpCode{Tag: 'r', I1: 0, I2: 3, J1: 0, J2: 1},
			verify: func(t *testing.T, result []DiffLine) {
				t.Helper()
				assert.Len(t, result, 3) // maxLen = 3
				// First line: both sides present
				assert.Equal(t, "line1", result[0].LeftLine)
				assert.Equal(t, "new1", result[0].RightLine)
				assert.NotEqual(t, 0, result[0].LeftNum)
				assert.NotEqual(t, 0, result[0].RightNum)

				// Second line: only left side (covers lines 112-116)
				assert.Equal(t, "line2", result[1].LeftLine)
				assert.Equal(t, "", result[1].RightLine)
				assert.NotEqual(t, 0, result[1].LeftNum)
				assert.Equal(t, 0, result[1].RightNum)

				// Third line: only left side
				assert.Equal(t, "line3", result[2].LeftLine)
				assert.Equal(t, "", result[2].RightLine)
			},
		},
		{
			name:     "more new lines than old",
			oldLines: []string{"old1"},
			newLines: []string{"new1", "new2", "new3"},
			opcode:   difflib.OpCode{Tag: 'r', I1: 0, I2: 1, J1: 0, J2: 3},
			verify: func(t *testing.T, result []DiffLine) {
				t.Helper()
				assert.Len(t, result, 3) // maxLen = 3
				// First line: both sides present
				assert.Equal(t, "old1", result[0].LeftLine)
				assert.Equal(t, "new1", result[0].RightLine)

				// Second line: only right side (covers lines 117-121)
				assert.Equal(t, "", result[1].LeftLine)
				assert.Equal(t, "new2", result[1].RightLine)
				assert.Equal(t, 0, result[1].LeftNum)
				assert.NotEqual(t, 0, result[1].RightNum)

				// Third line: only right side
				assert.Equal(t, "", result[2].LeftLine)
				assert.Equal(t, "new3", result[2].RightLine)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _, _ := processReplaceLines(tt.opcode, tt.oldLines, tt.newLines, 1, 1)
			tt.verify(t, result)
		})
	}
}

func Test_compressContext_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   []DiffLine
		context int
		verify  func(t *testing.T, result []DiffLine)
	}{
		{
			name: "no changes - all context",
			input: []DiffLine{
				{LineType: DiffLineContext, LeftLine: "line1", RightLine: "line1", LeftNum: 1, RightNum: 1},
				{LineType: DiffLineContext, LeftLine: "line2", RightLine: "line2", LeftNum: 2, RightNum: 2},
			},
			context: 2,
			verify: func(t *testing.T, result []DiffLine) {
				t.Helper()
				// Should return as-is since no changes (covers line 179-181)
				assert.Len(t, result, 2)
			},
		},
		{
			name: "change at start",
			input: []DiffLine{
				{LineType: DiffLineAdded, LeftLine: "", RightLine: "new", LeftNum: 0, RightNum: 1},
				{LineType: DiffLineContext, LeftLine: "ctx1", RightLine: "ctx1", LeftNum: 1, RightNum: 2},
				{LineType: DiffLineContext, LeftLine: "ctx2", RightLine: "ctx2", LeftNum: 2, RightNum: 3},
				{LineType: DiffLineContext, LeftLine: "ctx3", RightLine: "ctx3", LeftNum: 3, RightNum: 4},
				{LineType: DiffLineContext, LeftLine: "ctx4", RightLine: "ctx4", LeftNum: 4, RightNum: 5},
			},
			context: 1,
			verify: func(t *testing.T, result []DiffLine) {
				t.Helper()
				// Should keep change + 1 context after
				assert.LessOrEqual(t, len(result), 3) // change + 1 context + compression marker
			},
		},
		{
			name: "change at end",
			input: []DiffLine{
				{LineType: DiffLineContext, LeftLine: "ctx1", RightLine: "ctx1", LeftNum: 1, RightNum: 1},
				{LineType: DiffLineContext, LeftLine: "ctx2", RightLine: "ctx2", LeftNum: 2, RightNum: 2},
				{LineType: DiffLineContext, LeftLine: "ctx3", RightLine: "ctx3", LeftNum: 3, RightNum: 3},
				{LineType: DiffLineRemoved, LeftLine: "old", RightLine: "", LeftNum: 4, RightNum: 0},
			},
			context: 1,
			verify: func(t *testing.T, result []DiffLine) {
				t.Helper()
				// Should keep 1 context before + change
				assert.LessOrEqual(t, len(result), 3) // compression + context + change
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compressContext(tt.input, tt.context)
			tt.verify(t, result)
		})
	}
}

func Test_createDiffStyleFunc_BoundaryConditions(t *testing.T) {
	styler := testStyle
	diffLines := []DiffLine{
		{LineType: DiffLineContext, LeftLine: "ctx", RightLine: "ctx", LeftNum: 1, RightNum: 1},
		{LineType: DiffLineRemoved, LeftLine: "old", RightLine: "", LeftNum: 2, RightNum: 0},
		{LineType: DiffLineAdded, LeftLine: "", RightLine: "new", LeftNum: 0, RightNum: 2},
	}

	styleFunc := createDiffStyleFunc(diffLines, styler)

	tests := []struct {
		name string
		row  int
		col  int
	}{
		{name: "header row left", row: -1, col: 0},
		{name: "header row right", row: -1, col: 1},
		{name: "first data row left", row: 0, col: 0},
		{name: "first data row right", row: 0, col: 1},
		{name: "removed line left", row: 1, col: 0},
		{name: "removed line right", row: 1, col: 1},
		{name: "added line left", row: 2, col: 0},
		{name: "added line right", row: 2, col: 1},
		{name: "out of bounds row", row: 100, col: 0},
		{name: "invalid col", row: 0, col: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic for any valid or invalid combination
			style := styleFunc(tt.row, tt.col)
			assert.NotNil(t, style)
		})
	}
}
