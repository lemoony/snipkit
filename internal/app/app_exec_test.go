package app

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_App_Exec(t *testing.T) {
	snippetContent := `# some comment
# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the terminal first
echo "${VAR1}"`

	snippets := []model.Snippet{
		{
			UUID:         "uuid1",
			TitleFunc:    testutil.FixedString("title-1"),
			LanguageFunc: testutil.FixedLanguage(model.LanguageYAML),
			TagUUIDs:     []string{},
			ContentFunc:  testutil.FixedString(snippetContent),
		},
	}

	inputVar1Value := "foo-value"

	terminal := uiMocks.Terminal{}
	terminal.On("ApplyConfig", mock.Anything, mock.Anything).Return()
	terminal.On("ShowLookup", mock.Anything).Return(0)
	terminal.On("ShowParameterForm", mock.Anything, mock.Anything).Return([]string{inputVar1Value, ""}, true)

	terminal.On(mockutil.PrintMessage, inputVar1Value+"\n").Return()

	app := NewApp(
		WithTerminal(&terminal), WithConfig(configtest.NewTestConfig().Config), withProviderSnippets(snippets),
	)

	app.LookupAndExecuteSnippet()
}
