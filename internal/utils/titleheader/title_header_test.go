package titleheader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PruneTitleHeader(t *testing.T) {
	tests := []struct {
		name     string
		snippet  string
		expected string
	}{
		{name: "example1", snippet: "#\r\n# Hello\r\n#\n\r\nfoo", expected: "foo"},
		{name: "example2", snippet: "#\n# Get PIDs which listens to port\n#\n\n\nfoo content", expected: "foo content"},
		{name: "example3", snippet: "#\n# title\n#", expected: ""},
		{name: "example4", snippet: "#\n# title", expected: "#\n# title"},
		{name: "example5", snippet: "#/bin/bash\n#\n#title\n#", expected: "#/bin/bash"},
		{
			name:     "example5",
			snippet:  "#/bin/bash\n#\n#title\n#\n\n# ${VAR} Name: Variable\necho${VAR}",
			expected: "#/bin/bash\n\n# ${VAR} Name: Variable\necho${VAR}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, PruneTitleHeader(tt.snippet))
		})
	}
}

func Test_ParseTitleFromHeader(t *testing.T) {
	tests := []struct {
		title   string
		content string
		ok      bool
	}{
		{title: "title 1", content: "#\n# title 1\n#", ok: true},
		{title: "title 2", content: "#\n#title 2\n#", ok: true},
		{title: "title 3", content: "#/bin/bash\n#\n#title 3\n#", ok: true},
		{title: "title 4", content: "#/bin/bash\n\n#\n#title 4\n#", ok: true},
		{title: "title 5", content: "#\n#title 2", ok: false},
		{title: "title 6", content: "#title 2\n#", ok: false},
		{title: "title 7", content: "#\n# \n#", ok: false},
		{title: "title 8", content: "\n\n\n#\n# title 8\n#", ok: false},
		{title: "title 9", content: "\n\n#\n# title 9\n#", ok: true},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			name, ok := ParseTitleFromHeader(tt.content)
			assert.Equal(t, tt.ok, ok)
			if tt.ok {
				assert.Equal(t, tt.title, name)
			}
		})
	}
}
