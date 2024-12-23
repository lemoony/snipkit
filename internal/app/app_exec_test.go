package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_App_Exec(t *testing.T) {
	snippets := []model.Snippet{
		testutil.TestSnippet{
			ID:       "uuid1",
			Title:    "title-1",
			Language: model.LanguageYAML,
			Tags:     []string{},
			Content:  testSnippetContent,
		},
	}

	inputVar1Value := "foo-value"

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowLookup", mock.Anything, mock.Anything).Return(0)
	tui.On("ShowParameterForm", mock.Anything, mock.Anything, mock.Anything).Return([]string{inputVar1Value, ""}, true)
	tui.On(mockutil.Print, mock.Anything)
	tui.On(mockutil.Confirmation, mock.Anything).Return(true)
	tui.On(mockutil.PrintMessage, inputVar1Value+"\n").Return()

	app := NewApp(
		WithTUI(&tui),
		WithConfig(configtest.NewTestConfig().Config),
		withManagerSnippets(snippets),
	)

	app.LookupAndExecuteSnippet(true, true)

	// TODO fix
	// tui.AssertCalled(t, mockutil.Confirmation, uimsg.ExecConfirm("title-1", testSnippetContent))
	// tui.AssertCalled(t, mockutil.Print, uimsg.ExecPrint("title-1", testSnippetContent))
}

func Test_App_Exec_FindScriptAndExecuteWithParameters(t *testing.T) {
	snippetContent := `# some comment
# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the tui first
echo "${VAR1}"`

	snippets := []model.Snippet{
		testutil.TestSnippet{
			ID:       "uuid1",
			Title:    "title-1",
			Language: model.LanguageYAML,
			Tags:     []string{},
			Content:  snippetContent,
		},
	}

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()

	app := NewApp(
		WithTUI(&tui),
		WithConfig(configtest.NewTestConfig().Config),
		withManagerSnippets(snippets),
	)

	app.FindScriptAndExecuteWithParameters("uuid1", []model.ParameterValue{{Key: "VAR1", Value: "foo"}}, false, false)
}

func Test_App_Exec_FindScriptAndExecuteWithParameters_MissingParameters(t *testing.T) {
	snippetContent := `# some comment
# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the tui first
echo "${VAR1}"`

	snippets := []model.Snippet{
		testutil.TestSnippet{
			ID:       "uuid1",
			Title:    "title-1",
			Language: model.LanguageYAML,
			Tags:     []string{},
			Content:  snippetContent,
		},
	}

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowParameterForm", mock.Anything, mock.Anything, mock.Anything).Return([]string{"VAR1", ""}, true)

	app := NewApp(
		WithTUI(&tui),
		WithConfig(configtest.NewTestConfig().Config),
		withManagerSnippets(snippets),
	)

	app.FindScriptAndExecuteWithParameters("uuid1", []model.ParameterValue{}, false, false)
	tui.AssertCalled(t, "ShowParameterForm", snippets[0].GetParameters(), []model.ParameterValue{}, ui.OkButtonExecute)
}

func Test_formatOptions(t *testing.T) {
	tests := []struct {
		config   config.ScriptConfig
		expected model.SnippetFormatOptions
	}{
		{
			config:   config.ScriptConfig{RemoveComments: true, ParameterMode: config.ParameterModeSet},
			expected: model.SnippetFormatOptions{RemoveComments: true, ParamMode: model.SnippetParamModeSet},
		},
		{
			config:   config.ScriptConfig{RemoveComments: false, ParameterMode: config.ParameterModeReplace},
			expected: model.SnippetFormatOptions{RemoveComments: false, ParamMode: model.SnippetParamModeReplace},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			assert.Equal(t, tt.expected, formatOptions(tt.config))
		})
	}
}
