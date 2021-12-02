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
echo "2 -> ${VAR2}"
	`

	parameters := ParseParameters(rawSnippets)
	assert.Len(t, parameters, 2)

	assert.Equal(t, "VAR1", parameters[0].Key)
	assert.Equal(t, "First Output", parameters[0].Name)
	assert.Equal(t, "What to print on the terminal first", parameters[0].Description)

	assert.Equal(t, "VAR2", parameters[1].Key)
	assert.Equal(t, "Second Output", parameters[1].Name)
	assert.Equal(t, "What to print on the terminal second", parameters[1].Description)
}
