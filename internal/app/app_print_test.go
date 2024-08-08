package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_LookupAndCreatePrintableSnippet(t *testing.T) {
	snippetContent := `# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the tui first
echo "${VAR1}`

	snippets := []model.Snippet{
		testutil.TestSnippet{ID: "uuid1", Title: "title-1", Language: model.LanguageYAML, Tags: []string{}, Content: "content-1"},
		testutil.TestSnippet{ID: "uuid2", Title: "title-2", Language: model.LanguageBash, Tags: []string{}, Content: snippetContent},
	}

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowLookup", mock.Anything, mock.Anything).Return(1)
	tui.On("ShowParameterForm", mock.Anything, mock.Anything, mock.Anything).Return([]string{"foo-value"}, true)

	app := NewApp(
		WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManagerSnippets(snippets),
	)

	s, ok := app.LookupAndCreatePrintableSnippet()
	assert.True(t, ok)
	assert.Equal(t, `# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the tui first
VAR1="foo-value"
echo "${VAR1}`, s)
}
