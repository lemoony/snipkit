package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_LookupAndCreatePrintableSnippet(t *testing.T) {
	snippetContent := `# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the tui first
echo "${VAR1}`

	snippets := []model.Snippet{
		{UUID: "uuid1", TitleFunc: testutil.FixedString("title-1"), LanguageFunc: testutil.FixedLanguage(model.LanguageYAML), TagUUIDs: []string{}, ContentFunc: testutil.FixedString("content-1")},
		{UUID: "uuid2", TitleFunc: testutil.FixedString("title-2"), LanguageFunc: testutil.FixedLanguage(model.LanguageBash), TagUUIDs: []string{}, ContentFunc: testutil.FixedString(snippetContent)},
	}

	tui := uiMocks.TUI{}
	tui.On("ApplyConfig", mock.Anything, mock.Anything).Return()
	tui.On("ShowLookup", mock.Anything).Return(1)
	tui.On("ShowParameterForm", mock.Anything, mock.Anything).Return([]string{"foo-value"}, true)

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
