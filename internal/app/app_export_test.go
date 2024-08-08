package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_ExportSnippets(t *testing.T) {
	snippets := []model.Snippet{
		testutil.TestSnippet{ID: "uuid1", Title: "title-1", Language: model.LanguageYAML, Tags: []string{}, Content: "# ${VAR} Name: Message\necho ${VAR}"},
		testutil.TestSnippet{ID: "uuid2", Title: "title-2", Language: model.LanguageBash, Tags: []string{}, Content: "content-2"},
	}

	app := NewApp(WithConfig(configtest.NewTestConfig().Config), withManagerSnippets(snippets))

	tests := []struct {
		fields   []ExportField
		expected string
	}{
		{
			fields:   []ExportField{ExportFieldID, ExportFieldTitle, ExportFieldContent, ExportFieldParameters},
			expected: `{"snippets":[{"id":"uuid1","title":"title-1","content":"# ${VAR} Name: Message\necho ${VAR}","parameters":[{"key":"VAR","name":"Message","type":"VALUE"}]},{"id":"uuid2","title":"title-2","content":"content-2"}]}`,
		},
		{
			fields:   []ExportField{ExportFieldID},
			expected: `{"snippets":[{"id":"uuid1"},{"id":"uuid2"}]}`,
		},
		{
			fields:   []ExportField{ExportFieldTitle},
			expected: `{"snippets":[{"title":"title-1"},{"title":"title-2"}]}`,
		},
		{
			fields:   []ExportField{ExportFieldContent},
			expected: `{"snippets":[{"content":"# ${VAR} Name: Message\necho ${VAR}"},{"content":"content-2"}]}`,
		},
		{
			fields:   []ExportField{ExportFieldID, ExportFieldParameters},
			expected: `{"snippets":[{"id":"uuid1","parameters":[{"key":"VAR","name":"Message","type":"VALUE"}]},{"id":"uuid2"}]}`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			s := app.ExportSnippets(tt.fields, ExportFormatJSON)
			assert.Equal(t, tt.expected, s)
		})
	}
}

func Test_ExportSnippets_formats(t *testing.T) {
	snippets := []model.Snippet{
		testutil.TestSnippet{ID: "uuid1", Title: "title-1", Content: "content-1"},
		testutil.TestSnippet{ID: "uuid2", Title: "title-2", Content: "content-2"},
	}

	app := NewApp(WithConfig(configtest.NewTestConfig().Config), withManagerSnippets(snippets))

	tests := []struct {
		format   ExportFormat
		expected string
	}{
		{
			format:   ExportFormatJSON,
			expected: `{"snippets":[{"id":"uuid1","content":"content-1"},{"id":"uuid2","content":"content-2"}]}`,
		},
		{
			format: ExportFormatPrettyJSON,
			expected: `{
    "snippets": [
        {
            "id": "uuid1",
            "content": "content-1"
        },
        {
            "id": "uuid2",
            "content": "content-2"
        }
    ]
}`,
		},
		{
			format: ExportFormatXML,
			expected: `<exportJSON>
    <Snippets>
        <id>uuid1</id>
        <content>content-1</content>
    </Snippets>
    <Snippets>
        <id>uuid2</id>
        <content>content-2</content>
    </Snippets>
</exportJSON>`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			s := app.ExportSnippets([]ExportField{ExportFieldID, ExportFieldContent}, tt.format)
			assert.Equal(t, tt.expected, s)
		})
	}
}
