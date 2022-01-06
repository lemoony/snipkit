package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snippet-kit/internal/config/configtest"
	"github.com/lemoony/snippet-kit/internal/model"
	uiMocks "github.com/lemoony/snippet-kit/mocks/ui"
)

func Test_LookupAndCreatePrintableSnippet(t *testing.T) {
	snippetContent := `# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the terminal first
echo "${VAR1}`

	snippets := []model.Snippet{
		{UUID: "uuid1", Title: "title-1", Language: model.LanguageYAML, TagUUIDs: []string{}, Content: "content-1"},
		{UUID: "uuid2", Title: "title-2", Language: model.LanguageBash, TagUUIDs: []string{}, Content: snippetContent},
	}

	terminal := uiMocks.Terminal{}
	terminal.On("ApplyConfig", mock.Anything).Return()
	terminal.On("ShowLookup", snippets).Return(1)
	terminal.On("ShowParameterForm", mock.Anything).Return([]string{"foo-value"})

	app := NewApp(
		WithTerminal(&terminal), WithConfig(configtest.NewTestConfig().Config), withProviderSnippets(snippets),
	)

	s := app.LookupAndCreatePrintableSnippet()
	assert.Equal(t, `# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the terminal first
VAR1="foo-value"
echo "${VAR1}`, s)
}
