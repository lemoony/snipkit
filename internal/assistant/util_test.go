package assistant

import (
	"fmt"
	"strings"
	"testing"
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
