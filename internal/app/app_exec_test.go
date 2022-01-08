package app

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snippet-kit/internal/config/configtest"
	"github.com/lemoony/snippet-kit/internal/model"
	uiMocks "github.com/lemoony/snippet-kit/mocks/ui"
)

func Test_App_Exec(t *testing.T) {
	snippetContent := `# some comment
# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the terminal first
echo "${VAR1}"`

	snippets := []model.Snippet{
		{UUID: "uuid1", Title: "title-1", Language: model.LanguageYAML, TagUUIDs: []string{}, Content: snippetContent},
	}

	inputVar1Value := "foo-value"

	terminal := uiMocks.Terminal{}
	terminal.On("ApplyConfig", mock.Anything, mock.Anything).Return()
	terminal.On("ShowLookup", snippets).Return(0)
	terminal.On("ShowParameterForm", mock.Anything).Return([]string{inputVar1Value, ""})

	terminal.On("PrintMessage", inputVar1Value+"\n").Return()

	app := NewApp(
		WithTerminal(&terminal), WithConfig(configtest.NewTestConfig().Config), withProviderSnippets(snippets),
	)

	app.LookupAndExecuteSnippet()

	terminal.AssertCalled(t, "PrintMessage", inputVar1Value+"\n")
}
