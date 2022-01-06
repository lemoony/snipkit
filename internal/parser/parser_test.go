package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/model"
)

const (
	testSnippet1 = `
# some comment
# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the terminal first
echo "1 -> ${VAR1}"

# ${VAR2} Name: Second Output
# ${VAR2} Description: What to print on the terminal second
# ${VAR2} Default: Hey there!
echo "2 -> ${VAR2}"
`
	testSnippet2 = `
# ${VAR1} Name: Choose from
# ${VAR1} Description: What to print on the terminal
# ${VAR1} Values: One + some more, "Two",Three,  ,
# ${VAR1} Values: Four\, and some more, Five

echo "1 -> ${VAR1}"
`
	testSnippet3 = `
# ${VAR1} Description: What to print on the terminal
echo "1 -> ${VAR1}"
`
)

func Test_parseParameters(t *testing.T) {
	tests := []struct {
		name       string
		snippet    string
		parameters []model.Parameter
	}{
		{name: "multiple", snippet: testSnippet1, parameters: []model.Parameter{
			{Key: "VAR1", Name: "First Output", Description: "What to print on the terminal first"},
			{Key: "VAR2", Name: "Second Output", Description: "What to print on the terminal second", DefaultValue: "Hey there!"},
		}},
		{name: "single_enum", snippet: testSnippet2, parameters: []model.Parameter{
			{
				Key:         "VAR1",
				Name:        "Choose from",
				Description: "What to print on the terminal",
				Values:      []string{"One + some more", "\"Two\"", "Three", "Four, and some more", "Five"},
			},
		}},
		{name: "no_name", snippet: testSnippet3, parameters: []model.Parameter{
			{Key: "VAR1", Name: "VAR1", Description: "What to print on the terminal"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualParameters := ParseParameters(tt.snippet)
			assert.Len(t, actualParameters, len(tt.parameters))
			for i, expected := range tt.parameters {
				assert.Equal(t, expected, actualParameters[i])
			}
		})
	}
}

func Test_createSnippet(t *testing.T) {
	parameters := ParseParameters(testSnippet1)
	printable := CreateSnippet(testSnippet1, parameters, []string{"FOO-1", "FOO-2"})

	assert.Equal(t, `
# some comment
# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the terminal first
VAR1="FOO-1"
echo "1 -> ${VAR1}"

# ${VAR2} Name: Second Output
# ${VAR2} Description: What to print on the terminal second
# ${VAR2} Default: Hey there!
VAR2="FOO-2"
echo "2 -> ${VAR2}"
`, printable)
}
