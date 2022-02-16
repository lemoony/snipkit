package pet

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

const (
	testDataUserHome    = "testdata/userhome"
	testDataSnippetFile = testDataUserHome + "/.config/pet/snippet.toml"
)

func Test_parseSnippetFilePaths(t *testing.T) {
	s := testutil.NewTestSystem(system.WithUserHome(testDataUserHome))
	snippetFilePaths, err := parseSnippetFilePaths(s)
	assert.NoError(t, err)
	assert.Len(t, snippetFilePaths, 1)
	assert.Equal(t, testDataSnippetFile, snippetFilePaths[0])
}

func Test_parseSnippetFilePaths_notFound(t *testing.T) {
	s := testutil.NewTestSystem(system.WithUserHome("/foo/path"))
	snippetFilePaths, err := parseSnippetFilePaths(s)
	assert.NoError(t, err)
	assert.Empty(t, snippetFilePaths)
}

func Test_parseSnippetFilePaths_invalidConfigTOML(t *testing.T) {
	s := system.NewSystem(system.WithFS(afero.NewMemMapFs()))
	s.CreatePath(filepath.Join(s.UserHome(), defaultConfigPath))
	s.WriteFile(filepath.Join(s.UserHome(), defaultConfigPath), []byte("foo: 1"))

	result, err := parseSnippetFilePaths(s)
	assert.Error(t, err)
	assert.Empty(t, result)
}

func Test_parseSnippetsFromTOML(t *testing.T) {
	system := testutil.NewTestSystem()
	contents := string(system.ReadFile(testDataSnippetFile))

	snippets := parseSnippetsFromTOML(contents)
	assert.Len(t, snippets, 2)
	assert.NotEmpty(t, snippets[0].GetID())
	assert.Equal(t, "Echo something", snippets[0].GetTitle())
	assert.Equal(t,
		"echo <VAR1> && <VAR2=Snipkit> <VAR3=is a snippet manager for the terminal!>",
		snippets[0].GetContent(),
	)
	assert.Equal(t, model.LanguageBash, snippets[0].GetLanguage())
	assert.Len(t, snippets[0].GetParameters(), 3)
	assert.Len(t, snippets[0].GetTags(), 0)
	assert.Equal(t, snippets[0].Format([]string{"one", "two", "three"}, model.SnippetFormatOptions{}), "echo one && two three")

	assert.Equal(t, "Watches Kubernetes pods with refresh", snippets[1].GetTitle())
	assert.Equal(t, "watch -n 5 'kubectl get pods | grep <pattern>'", snippets[1].GetContent())
	assert.Equal(t, model.LanguageBash, snippets[1].GetLanguage())
	assert.Len(t, snippets[1].GetParameters(), 1)
	assert.Len(t, snippets[1].GetTags(), 2)
	assert.Equal(t, snippets[1].Format([]string{"foo"}, model.SnippetFormatOptions{}), "watch -n 5 'kubectl get pods | grep foo'")
}

func Test_parseParameters(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected []model.Parameter
	}{
		{name: "no parameters", command: `echo "hello world"`, expected: []model.Parameter{}},
		{name: "empty default", command: `echo "<var1=>"`, expected: []model.Parameter{{Key: "var1"}}},
		{
			name:    "multiple parameters with no default value",
			command: `echo "<var1>" && echo "<var2>"`,
			expected: []model.Parameter{
				{Key: "var1"},
				{Key: "var2"},
			},
		},
		{
			name:    "parameter with value",
			command: `echo "<var1=foo>"`,
			expected: []model.Parameter{
				{Key: "var1", DefaultValue: "foo"},
			},
		},
		{
			name:    "multiple with value",
			command: `kubectl config use-context k8s-<hub=emea>-<environment=e2e>`,
			expected: []model.Parameter{
				{Key: "hub", DefaultValue: "emea"},
				{Key: "environment", DefaultValue: "e2e"},
			},
		},
		{
			name:    "value with whitespace",
			command: `echo <var1> && <var2=Snipkit> <var3=is a snippet manager for the terminal!>`,
			expected: []model.Parameter{
				{Key: "var1"},
				{Key: "var2", DefaultValue: "Snipkit"},
				{Key: "var3", DefaultValue: "is a snippet manager for the terminal!"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualParameters := parseParameters(tt.command)
			if len(tt.expected) == 0 {
				assert.Empty(t, actualParameters)
			} else {
				assert.Equal(t, tt.expected, actualParameters)
			}
		})
	}
}

func Test_formatContent(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		values   []string
		expected string
	}{
		{name: "no parameters", command: `echo "hello world"`, expected: `echo "hello world"`},
		{name: "empty default and value", command: `echo "<var1=>"`, values: []string{""}, expected: `echo ""`},
		{name: "parameter with value", command: `echo "<var1=foo>"`, values: []string{"test"}, expected: `echo "test"`},
		{
			name:     "value with whitespace",
			command:  `echo <var1> <var2=Snipkit> && echo <var3=is a snippet manager for the terminal!>`,
			values:   []string{"one", "two", "three"},
			expected: "echo one two && echo three",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, formatContent(tt.command, tt.values))
		})
	}
}
