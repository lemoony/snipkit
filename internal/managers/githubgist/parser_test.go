package githubgist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
)

func Test_parseSnippet(t *testing.T) {
	raw := rawSnippet{
		ID: "some-id",
		Content: []byte(`#
# Snippet Title
#
echo "Hello World"`),
		Language:    "Shell",
		Pubic:       true,
		FilesInGist: 1,
		Filename:    "echo-something.sh",
		ETag:        "etag",
		Description: "Echo something #test",
	}

	cfg := GistConfig{
		HideTitleInPreview: true,
		TitleHeaderEnabled: true,
		NameMode:           SnippetNameModeDescription,
	}

	snippet := parseSnippet(raw, cfg)

	assert.Equal(t, "Snippet Title", snippet.GetTitle())
	assert.Equal(t, `echo "Hello World"`, snippet.GetContent())
	assert.Equal(t, model.LanguageBash, snippet.GetLanguage())
	assert.Equal(t, []string{"test"}, snippet.TagUUIDs)
}

func Test_parseTitle(t *testing.T) {
	type testTemplate struct {
		description        string
		filename           string
		filesInGist        int
		content            string
		nameMode           SnippetNameMode
		titleHeaderEnabled bool
		expected           string
	}

	testWith := func(nameMode SnippetNameMode, titleHeaderEnabled bool, filesInGist int, expected string) testTemplate {
		return testTemplate{
			description:        "snippet description",
			filename:           "filename.sh",
			filesInGist:        filesInGist,
			content:            "",
			nameMode:           nameMode,
			titleHeaderEnabled: titleHeaderEnabled,
			expected:           expected,
		}
	}

	testWithContent := func(content, expected string) testTemplate {
		test := testWith(SnippetNameModeDescription, true, 1, expected)
		test.content = content
		return test
	}

	tests := []testTemplate{
		testWith(SnippetNameModeDescription, true, 1, "snippet description"),
		testWith(SnippetNameModeFilename, true, 1, "filename.sh"),
		testWith(SnippetNameModeCombine, true, 1, "snippet description - filename.sh"),
		testWith(SnippetNameModeCombinePreferDescription, true, 1, "snippet description"),
		testWith(SnippetNameModeCombinePreferFilename, true, 1, "filename.sh"),
		testWith(SnippetNameModeCombinePreferDescription, true, 2, "snippet description - filename.sh"),
		testWith(SnippetNameModeCombinePreferFilename, true, 2, "snippet description - filename.sh"),
		testWithContent("#\n# Title Header\n#", "Title Header"),
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			title := parseTitle(
				rawSnippet{
					ID:          "foo-id",
					ETag:        "",
					Description: tt.description,
					Filename:    tt.filename,
					FilesInGist: tt.filesInGist,
					Content:     []byte(tt.content),
					Pubic:       true,
				},
				tt.nameMode,
				tt.titleHeaderEnabled,
			)

			assert.Equal(t, tt.expected, title)
		})
	}
}

func Test_parseTags(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{name: "single tag", text: "#foo", expected: []string{"foo"}},
		{name: "one tag with hash", text: "#foo#hey", expected: []string{"foo#hey"}},
		{name: "no tag", text: "foo", expected: []string{}},
		{name: "text with multiple tags", text: "Hello world #foo #stuff", expected: []string{"foo", "stuff"}},
		{name: "empty text", text: "", expected: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := parseTags(tt.text)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_pruneTags(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{name: "empty text", text: "", expected: ""},
		{name: "one tag with hash", text: "#foo#hey", expected: ""},
		{name: "text with no tag", text: "foo", expected: "foo"},
		{name: "text with multiple tags", text: "Hello world #foo #stuff", expected: "Hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, pruneTags(tt.text))
		})
	}
}

func Test_formatContent(t *testing.T) {
	tests := []struct {
		name            string
		value           string
		hideTitleHeader bool
		expected        string
	}{
		{name: "no title header", value: "Content", hideTitleHeader: true, expected: "Content"},
		{name: "title header", value: "#\n# Title\n#\nContent", hideTitleHeader: true, expected: "Content"},
		{name: "don't hide title", value: "#\n# Title\n#\nContent", hideTitleHeader: false, expected: "#\n# Title\n#\nContent"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, formatContent(tt.value, tt.hideTitleHeader))
		})
	}
}

func Test_mapLanguage(t *testing.T) {
	tests := []struct {
		value    string
		expected model.Language
	}{
		{value: "YAML", expected: model.LanguageYAML},
		{value: "Shell", expected: model.LanguageBash},
		{value: "TOML", expected: model.LanguageTOML},
		{value: "Markdown", expected: model.LanguageMarkdown},
		{value: "Foo", expected: model.LanguageText},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equal(t, tt.expected, mapLanguage(tt.value))
		})
	}
}
