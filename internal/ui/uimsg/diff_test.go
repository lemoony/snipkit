package uimsg

import (
	"strings"
	"testing"

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
