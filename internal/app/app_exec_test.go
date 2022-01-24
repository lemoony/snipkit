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
# ${VAR1} Description: What to print on the tui first
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

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowLookup", mock.Anything).Return(0)
	tui.On("ShowParameterForm", mock.Anything, mock.Anything).Return([]string{inputVar1Value, ""}, true)

	tui.On(mockutil.PrintMessage, inputVar1Value+"\n").Return()

	app := NewApp(
		WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManagerSnippets(snippets),
	)

	app.LookupAndExecuteSnippet()
}
