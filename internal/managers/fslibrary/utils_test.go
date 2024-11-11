package fslibrary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_formatSnippet(t *testing.T) {
	tests := []struct {
		name     string
		script   string
		title    string
		expected string
	}{
		{
			name:   "Script with shebang",
			script: "#!/bin/bash\necho \"Hello, World!\"",
			title:  "Hello World Script",
			expected: `#!/bin/bash

#
# Hello World Script
#

echo "Hello, World!"`,
		},
		{
			name:   "Script without shebang",
			script: "echo \"Hello, World!\"",
			title:  "Hello World Script",
			expected: `#
# Hello World Script
#

echo "Hello, World!"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := formatSnippet(tt.script, tt.title)
			assert.Equal(t, tt.expected, formatted)
		})
	}
}
