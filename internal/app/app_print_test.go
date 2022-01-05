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
	snippetContent := `# some comment
# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the terminal first
echo "${VAR1}
# ${VAR2} Name: Second Output
# ${VAR2} Description: What to print on the terminal second
# ${VAR2} Default: default-value
echo "${VAR2}"
`

	snippets := []model.Snippet{
		{UUID: "uuid1", Title: "title-1", Language: model.LanguageYAML, TagUUIDs: []string{}, Content: "content-1"},
		{UUID: "uuid2", Title: "title-2", Language: model.LanguageBash, TagUUIDs: []string{}, Content: snippetContent},
	}

	terminal := uiMocks.Terminal{}
	terminal.On("ApplyConfig", mock.Anything).Return()
	terminal.On("ShowLookup", snippets).Return(1)
	terminal.On("ShowParameterForm", mock.Anything).Return([]string{"foo-value", ""})

	app := NewApp(
		WithTerminal(&terminal), WithConfig(configtest.NewTestConfig().Config), withProviderSnippets(snippets),
	)

	s := app.LookupAndCreatePrintableSnippet()
	assert.Equal(t, `# some comment
# foo-value Name: First Output
# foo-value Description: What to print on the terminal first
echo "foo-value
# default-value Name: Second Output
# default-value Description: What to print on the terminal second
# default-value Default: default-value
echo "default-value"
`, s)
}
