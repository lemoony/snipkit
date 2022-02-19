package form

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_suggestionsForPath(t *testing.T) {
	fs := afero.NewMemMapFs()

	createTestFile(t, fs, "testfile-a.txt")
	createTestFile(t, fs, "testfile-b.txt")

	createTestDirectory(t, fs, ".config/logs")
	createTestFile(t, fs, ".config/some.txt")

	createTestFile(t, fs, ".config/logs/file.log")
	createTestFile(t, fs, ".config/logs/HEAD")

	tests := []struct {
		path     string
		expected []string
	}{
		{path: "test", expected: []string{"testfile-a.txt", "testfile-b.txt"}},
		{path: "testfile-a.txt", expected: []string{}},
		{path: "./", expected: []string{"./", "./.config", "./testfile-a.txt", "./testfile-b.txt"}},
		{path: ".", expected: []string{".", "./.config", "./testfile-a.txt", "./testfile-b.txt"}},
		{path: "./test", expected: []string{"./testfile-a.txt", "./testfile-b.txt"}},
		{path: "./.config", expected: []string{"./.config", "./.config/logs", "./.config/some.txt"}},
		{path: ".config", expected: []string{".config", ".config/logs", ".config/some.txt"}},
		{path: ".config/logs", expected: []string{".config/logs", ".config/logs/HEAD", ".config/logs/file.log"}},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if len(tt.expected) == 0 {
				assert.Empty(t, suggestionsForPath(fs, tt.path))
			} else {
				assert.Equal(t, tt.expected, suggestionsForPath(fs, tt.path))
			}
		})
	}
}

func createTestFile(t *testing.T, fs afero.Fs, path string) {
	t.Helper()
	const fileMode = 0o600
	assert.NoError(t, afero.WriteFile(fs, path, []byte("foo"), fileMode))
}

func createTestDirectory(t *testing.T, fs afero.Fs, path string) {
	t.Helper()
	const dirMode = 0o700
	assert.NoError(t, fs.MkdirAll(path, dirMode))
}
