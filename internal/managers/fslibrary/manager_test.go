package fslibrary

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_GetInfo(t *testing.T) {
	system := testutil.NewTestSystem()
	libraryPath := t.TempDir()
	config := Config{
		Enabled:     true,
		LibraryPath: []string{libraryPath},
		SuffixRegex: []string{".sh", ".yaml"},
	}

	provider, err := NewManager(WithSystem(system), WithConfig(config))
	assert.NoError(t, err)

	info := provider.Info()

	assert.Len(t, info.Lines, 3)

	assert.Equal(t, info.Lines[0].Key, "Filesystem library paths")
	assert.Equal(t, info.Lines[0].Value, fmt.Sprintf("[%s]", libraryPath))
	assert.False(t, info.Lines[0].IsError)

	assert.Equal(t, info.Lines[1].Key, "Filesystem library allowed suffixes")
	assert.Equal(t, info.Lines[1].Value, "[.sh, .yaml]")
	assert.False(t, info.Lines[1].IsError)

	assert.Equal(t, info.Lines[2].Key, "Filesystem library total number of snippets")
	assert.Equal(t, info.Lines[2].Value, "0")
	assert.False(t, info.Lines[2].IsError)
}

func Test_GetSnippets(t *testing.T) {
	system := testutil.NewTestSystem()
	config := Config{
		Enabled:     true,
		LibraryPath: []string{t.TempDir()},
		SuffixRegex: []string{".sh", ".yaml"},
	}

	const filePerm = 0o600

	files := []struct {
		file     string
		language model.Language
	}{
		{file: "snippet-0.sh", language: model.LanguageBash},
		{file: "snippet-1.sh", language: model.LanguageBash},
		{file: "snippet-2.yaml", language: model.LanguageYAML},
	}

	for i := 0; i < len(files); i++ {
		assert.NoError(t, afero.WriteFile(
			system.Fs,
			filepath.Join(config.LibraryPath[0], files[i].file),
			[]byte(fmt.Sprintf("content-%d", i)),
			filePerm,
		))
	}

	// write one file into library dir which does not match the suffix regex
	assert.NoError(t, afero.WriteFile(
		system.Fs,
		filepath.Join(config.LibraryPath[0], "foo.toml"),
		[]byte("foo"),
		filePerm,
	))

	provider, err := NewManager(WithSystem(system), WithConfig(config))
	assert.NoError(t, err)

	snippets := provider.GetSnippets()
	assert.Len(t, snippets, len(files))

	for i, s := range snippets {
		assert.Equal(t, files[i].file, s.GetTitle())
		assert.Equal(t, files[i].language, s.GetLanguage())
		assert.Equal(t, fmt.Sprintf("content-%d", i), s.GetContent())
	}
}

func Test_GetSnippets_LazyOpen_HideTitleHeader(t *testing.T) {
	system := testutil.NewTestSystem()
	config := Config{
		Enabled:            true,
		LazyOpen:           true,
		HideTitleInPreview: true,
		LibraryPath:        []string{t.TempDir()},
		SuffixRegex:        []string{".sh", ".yaml"},
	}

	const filePerm = 0o600

	assert.NoError(t, afero.WriteFile(
		system.Fs,
		filepath.Join(config.LibraryPath[0], "snippet.sh"),
		[]byte("#\n# title\n#\ncontent"),
		filePerm,
	))

	provider, err := NewManager(WithSystem(system), WithConfig(config))
	assert.NoError(t, err)

	snippets := provider.GetSnippets()
	assert.Len(t, snippets, 1)

	assert.Equal(t, "content", snippets[0].GetContent())
	assert.Equal(t, "snippet.sh", snippets[0].GetTitle())
}

func Test_checkSuffix(t *testing.T) {
	tests := []struct {
		filename string
		re       []*regexp.Regexp
		expected bool
	}{
		{filename: ".sh", re: []*regexp.Regexp{regexp.MustCompile(".sh")}, expected: true},
		{filename: ".yaml", re: []*regexp.Regexp{regexp.MustCompile(".sh")}, expected: false},
		{filename: ".sh", re: []*regexp.Regexp{}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			assert.Equal(t, tt.expected, checkSuffix(tt.filename, tt.re))
		})
	}
}

func Test_languageForSuffix(t *testing.T) {
	tests := []struct {
		suffix   string
		expected model.Language
	}{
		{suffix: ".sh", expected: model.LanguageBash},
		{suffix: ".yaml", expected: model.LanguageYAML},
		{suffix: ".yml", expected: model.LanguageYAML},
		{suffix: ".md", expected: model.LanguageMarkdown},
		{suffix: ".toml", expected: model.LanguageTOML},
		{suffix: ".txt", expected: model.LanguageUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.suffix, func(t *testing.T) {
			assert.Equal(t, tt.expected, languageForSuffix(tt.suffix))
		})
	}
}
