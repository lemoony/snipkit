package assistant

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
)

func Test_extractBashScript(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedScript   string
		expectedFilename string
	}{
		{
			name: "with markdown + filename",
			input: `#!/bin/sh
#
# Simple script
# Filename: simple-script.sh
#
echo "foo"`,
			expectedFilename: "simple-script.sh",
			expectedScript: `#!/bin/sh
#
# Simple script
#
echo "foo"`,
		},
		{
			name: "without markdown + no filename",
			input: wrapInMarkdown(`#!/bin/sh
echo "foo"`),
			expectedFilename: "",
			expectedScript: `#!/bin/sh
echo "foo"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script, filename := extractBashScript(tt.input)
			if strings.TrimSpace(script) != tt.expectedScript {
				t.Errorf("extractBashScript() got = %v, want %v", script, tt.expectedScript)
			}
			if filename != tt.expectedFilename {
				t.Errorf("extractBashScript() got1 = %v, want %v", filename, tt.expectedFilename)
			}
		})
	}
}

func wrapInMarkdown(input string) string {
	return fmt.Sprintf("```sh\n%s\n```", input)
}

func Test_RandomScriptFilename(t *testing.T) {
	assert.NotEmpty(t, RandomScriptFilename())
}

func TestPrepareSnippet(t *testing.T) {
	content := `#!/bin/sh
#
# Simple script
#
# ${FOO} Name: Foo
echo ${FOO}`

	snippet := PrepareSnippet([]byte(content))

	assert.Equal(t, "Simple script", snippet.GetTitle())
	assert.Equal(t, model.LanguageBash, snippet.GetLanguage())
	assert.Equal(t, content, snippet.GetContent())
	assert.Len(t, snippet.GetParameters(), 1)
	assert.Equal(t, "Foo", snippet.GetParameters()[0].Name)
}
