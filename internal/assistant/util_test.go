package assistant

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
)

//nolint:funlen // test function for yaml config is allowed to be too long
func Test_extractBashScript(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedScript   string
		expectedFilename string
		expectedTittle   string
	}{
		{
			name: "with markdown + Filename + Title",
			input: `#!/bin/sh
#
# Snippet Title: Simple Contents
# Filename: simple-Contents.sh
#
echo "foo"`,
			expectedFilename: "simple-Contents.sh",
			expectedTittle:   "Simple Contents",
			expectedScript: `#!/bin/sh
echo "foo"`,
		},
		{
			name: "with markdown + Filename",
			input: `#!/bin/sh
#
# Filename: simple-Contents.sh
#
echo "foo"`,
			expectedFilename: "simple-Contents.sh",
			expectedScript: `#!/bin/sh
echo "foo"`,
		},
		{
			name: "with markdown + Title",
			input: `#!/bin/sh
#
# Snippet Title: Simple Contents
#
echo "foo"`,
			expectedTittle: "Simple Contents",
			expectedScript: `#!/bin/sh
echo "foo"`,
		},
		{
			name: "without markdown + no Filename",
			input: wrapInMarkdown(`#!/bin/sh
echo "foo"`),
			expectedFilename: "",
			expectedTittle:   "",
			expectedScript: `#!/bin/sh
echo "foo"`,
		},
		{
			name: "without markdown + random comments",
			input: wrapInMarkdown(`#!/bin/sh
#
# Foo
#
echo "foo"`),
			expectedFilename: "",
			expectedTittle:   "",
			expectedScript: `#!/bin/sh
#
# Foo
#
echo "foo"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := parseScript(tt.input)
			if strings.TrimSpace(parsed.Contents) != tt.expectedScript {
				t.Errorf("parseScript() got = %v, want %v", parsed.Contents, tt.expectedScript)
			}
			if parsed.Filename != tt.expectedFilename {
				t.Errorf("parseScript() got1 = %v, want %v", parsed.Filename, tt.expectedFilename)
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
# Simple Contents
#
# ${FOO} Name: Foo
echo ${FOO}`

	snippet := PrepareSnippet([]byte(content))

	assert.Equal(t, "Simple Contents", snippet.GetTitle())
	assert.Equal(t, model.LanguageBash, snippet.GetLanguage())
	assert.Equal(t, content, snippet.GetContent())
	assert.Len(t, snippet.GetParameters(), 1)
	assert.Equal(t, "Foo", snippet.GetParameters()[0].Name)
}
