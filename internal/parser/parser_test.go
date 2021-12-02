package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseSimpleParameters(t *testing.T) {
	rawSnippets := `
# some comment
# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the terminal first
echo "1 -> ${VAR1}"

# ${VAR2} Name: Second Output
# ${VAR2} Description: What to print on the terminal second
# ${VAR2} Default: Hey there!
echo "2 -> ${VAR2}"
	`

	parameters := ParseParameters(rawSnippets)
	assert.Len(t, parameters, 2)

	assert.Equal(t, "VAR1", parameters[0].Key)
	assert.Equal(t, "First Output", parameters[0].Name)
	assert.Equal(t, "What to print on the terminal first", parameters[0].Description)
	assert.Empty(t, parameters[0].DefaultValue)

	assert.Equal(t, "VAR2", parameters[1].Key)
	assert.Equal(t, "Second Output", parameters[1].Name)
	assert.Equal(t, "What to print on the terminal second", parameters[1].Description)
	assert.Equal(t, "Hey there!", parameters[1].DefaultValue)
}

func Test_parseEnum(t *testing.T) {
	rawSnippets := `
# ${VAR1} Name: Choose from
# ${VAR1} Description: What to print on the terminal
# ${VAR1} Values: One + some more, "Two",Three,  ,
# ${VAR1} Values: Four\, and some more, Five

echo "1 -> ${VAR1}"
	`

	parameters := ParseParameters(rawSnippets)
	assert.Len(t, parameters, 1)

	assert.Equal(t, "VAR1", parameters[0].Key)
	assert.Equal(t, "Choose from", parameters[0].Name)
	assert.Equal(t, "What to print on the terminal", parameters[0].Description)
	assert.Empty(t, parameters[0].DefaultValue)

	values := parameters[0].Values
	assert.Len(t, values, 5)
	assert.Equal(t, "One + some more", values[0])
	assert.Equal(t, "\"Two\"", values[1])
	assert.Equal(t, "Three", values[2])
	assert.Equal(t, "Four, and some more", values[3])
	assert.Equal(t, "Five", values[4])
}
