package stringutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringOrDefault(t *testing.T) {
	defaultValue := "default-value"
	noDefault := "no-default"
	assert.Equal(t, defaultValue, StringOrDefault("", defaultValue))
	assert.Equal(t, noDefault, StringOrDefault("no-default", noDefault))
}

func Test_SplitWithEscape(t *testing.T) {
	s := "One, Two\\, plus more, Three,, five"

	splits := SplitWithEscape(s, ',', '\\', true)
	assert.Len(t, splits, 4)
	assert.Equal(t, "One", splits[0])
	assert.Equal(t, "Two, plus more", splits[1])
	assert.Equal(t, "Three", splits[2])
	assert.Equal(t, "five", splits[3])
}

func Test_SplitWithEscapeNoSplit(t *testing.T) {
	s := "String without split character"
	splits := SplitWithEscape(s, ',', '\\', true)
	assert.Len(t, splits, 1)
	assert.Equal(t, s, splits[0])
}

func Test_SplitWithSingleEscapedCharacter(t *testing.T) {
	s := "String without split\\, character"
	splits := SplitWithEscape(s, ',', '\\', true)
	assert.Len(t, splits, 1)
	assert.Equal(t, "String without split, character", splits[0])
}

func Test_SplitWithoutTrimming(t *testing.T) {
	s := " One, Two"
	splits := SplitWithEscape(s, ',', '\\', false)
	assert.Len(t, splits, 2)
	assert.Equal(t, " One", splits[0])
	assert.Equal(t, " Two", splits[1])
}

func Test_FirstNotEmpty(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected string
	}{
		{name: "no values", values: []string{}, expected: ""},
		{name: "nil", values: nil, expected: ""},
		{name: "first value", values: []string{"first", ""}, expected: "first"},
		{name: "last value", values: []string{"", "last"}, expected: "last"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, FirstNotEmpty(tt.values...))
		})
	}
}
