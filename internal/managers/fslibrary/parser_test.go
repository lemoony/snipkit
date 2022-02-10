package fslibrary

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_pruneTitleComment(t *testing.T) {
	tests := []struct {
		name     string
		snippet  string
		expected string
	}{
		{name: "example1", snippet: "#\r\n# Hello\r\n#\n\r\nfoo", expected: "foo"},
		{name: "example2", snippet: "#/bin/bash\n#\n#title\n#", expected: "#/bin/bash"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, pruneTitleHeader(strings.NewReader(tt.snippet)))
		})
	}
}

func Test_getSnippetName(t *testing.T) {
	tests := []struct {
		title   string
		content string
		ok      bool
	}{
		{title: "title 1", content: "#\n# title 1\n#", ok: true},
		{title: "title 2", content: "#title 2\n#", ok: false},
	}

	system := testutil.NewTestSystem()

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			file, err := afero.TempFile(system.Fs, t.TempDir(), "*.sh")

			assert.NoError(t, err)
			if _, err := file.Write([]byte(tt.content)); err != nil {
				assert.NoError(t, err)
			}

			name := getSnippetName(system, file.Name())
			if tt.ok {
				assert.Equal(t, tt.title, name)
			} else {
				assert.Equal(t, filepath.Base(file.Name()), name)
			}
		})
	}
}
