package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/ui/assistant/chat"
	"github.com/lemoony/snipkit/internal/ui/picker"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	assistantMocks "github.com/lemoony/snipkit/mocks/assistant"
	configMocks "github.com/lemoony/snipkit/mocks/config"
	managerMocks "github.com/lemoony/snipkit/mocks/managers"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_App_GenerateSnippetWithAssistant_SaveExit(t *testing.T) {
	const exampleFile = "echo-foo.sh"
	const exampleTitle = "Echo foo!"
	const exampleScript = `
#!/bin/bash
# ${PARAM} Key: FOO_KEY
echo ${FOO_KEY}
`

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	// First call: initial generation, user chooses Execute
	tui.On(mockutil.ShowUnifiedAssistantChat, mock.Anything).Return(
		assistant.ParsedScript{Contents: exampleScript, Filename: exampleFile, Title: exampleTitle},
		[]string{"hello world"}, // parameterValues - provide values so execution happens
		chat.PreviewActionExecute,
		"", // latestPrompt
		"", // saveFilename
		"", // saveSnippetName
	).Once()
	// Second call: after execution, user chooses Cancel with save data
	tui.On(mockutil.ShowUnifiedAssistantChat, mock.Anything).Return(
		assistant.ParsedScript{Contents: exampleScript, Filename: exampleFile, Title: exampleTitle},
		[]string{},               // parameterValues
		chat.PreviewActionCancel, // Exit after execution
		"",                       // latestPrompt
		exampleFile,              // saveFilename
		exampleTitle,             // saveSnippetName
	).Once()

	cfg := configtest.NewTestConfig().Config
	cfgService := configMocks.ConfigService{}
	cfgService.On("LoadConfig").Return(cfg, nil)
	cfgService.On("NeedsMigration").Return(false, "")

	assistantMock := assistantMocks.Assistant{}
	assistantMock.On("Query", mock.Anything).Return(assistant.ParsedScript{
		Contents: exampleScript, Filename: exampleFile, Title: exampleTitle,
	})
	assistantMock.On("Initialize").Return(true, uimsg.Printable{})

	fsLibManager := managerMocks.Manager{}
	fsLibManager.On("Key").Return(fslibrary.Key)
	fsLibManager.On(mockutil.SaveAssistantSnippet, mock.Anything, mock.Anything, mock.Anything).Return("/path", exampleFile)

	provider := managerMocks.Provider{}
	provider.On("CreateManager", mock.Anything, mock.Anything, mock.Anything, &tui).Return([]managers.Manager{&fsLibManager}, nil)

	app := NewApp(
		WithTUI(&tui),
		WithConfigService(&cfgService),
		WithProvider(&provider),
		WithAssistantProviderFunc(func(config assistant.Config, demoConfig assistant.DemoConfig) assistant.Assistant {
			return &assistantMock
		}),
	)

	app.GenerateSnippetWithAssistant([]string{}, 0)

	fsLibManager.AssertCalled(t, mockutil.SaveAssistantSnippet, exampleTitle, exampleFile, []byte(exampleScript))
}

func Test_App_GenerateSnippetWithAssistant_TweakPrompt_DontSave(t *testing.T) {
	const prompt1 = "prompt 1"
	const prompt2 = "The result of the command was: \n\n\nprompt 2"
	const exampleFile1 = "echo-foo-1.sh"
	const exampleFile2 = "echo-foo-2.sh"
	const exampleScript1 = `#!/bin/bash echo one`
	const exampleScript2 = `#!/bin/bash echo two`

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	// First script: generation, then execute
	tui.On(mockutil.ShowUnifiedAssistantChat, mock.Anything).Return(
		assistant.ParsedScript{Contents: exampleScript1, Filename: exampleFile1},
		[]string{}, // No params needed for these scripts
		chat.PreviewActionExecute,
		"", // latestPrompt
		"", // saveFilename
		"", // saveSnippetName
	).Once()
	// After first execution: user chooses Revise to try again
	tui.On(mockutil.ShowUnifiedAssistantChat, mock.Anything).Return(
		assistant.ParsedScript{Contents: exampleScript1, Filename: exampleFile1},
		[]string{},
		chat.PreviewActionRevise, // Try again
		prompt2,                  // latestPrompt - user enters new prompt
		"",                       // saveFilename
		"",                       // saveSnippetName
	).Once()
	// Second script: generation, then execute
	tui.On(mockutil.ShowUnifiedAssistantChat, mock.Anything).Return(
		assistant.ParsedScript{Contents: exampleScript2, Filename: exampleFile2},
		[]string{},
		chat.PreviewActionExecute,
		"", // latestPrompt
		"", // saveFilename
		"", // saveSnippetName
	).Once()
	// After second execution: user chooses ExitNoSave
	tui.On(mockutil.ShowUnifiedAssistantChat, mock.Anything).Return(
		assistant.ParsedScript{Contents: exampleScript2, Filename: exampleFile2},
		[]string{},
		chat.PreviewActionExitNoSave, // Don't save and exit
		"",                           // latestPrompt
		"",                           // saveFilename
		"",                           // saveSnippetName
	).Once()
	// No OpenEditor call since PreviewActionExecute skips the editor
	tui.On(mockutil.Confirmation, mock.Anything).Return(true)

	cfg := configtest.NewTestConfig().Config
	cfgService := configMocks.ConfigService{}
	cfgService.On("LoadConfig").Return(cfg, nil)
	cfgService.On("NeedsMigration").Return(false, "")

	assistantMock := assistantMocks.Assistant{}
	assistantMock.On(mockutil.Query, prompt1).Return(assistant.ParsedScript{Contents: exampleScript1, Filename: exampleFile1})
	assistantMock.On(mockutil.Query, mock.Anything).Return(assistant.ParsedScript{Contents: exampleScript2, Filename: exampleFile2})
	assistantMock.On(mockutil.ValidateConfig).Return(true, uimsg.Printable{})

	app := NewApp(
		WithTUI(&tui),
		WithConfigService(&cfgService),
		WithAssistantProviderFunc(func(config assistant.Config, demoConfig assistant.DemoConfig) assistant.Assistant {
			return &assistantMock
		}),
	)

	app.GenerateSnippetWithAssistant([]string{}, 0)
}

func Test_App_EnableAssistant(t *testing.T) {
	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On(mockutil.ShowPicker, mock.Anything, mock.Anything, mock.Anything).Return(1, true)
	tui.On(mockutil.Confirmation, mock.Anything).Return(true)
	tui.On(mockutil.Print, mock.Anything)

	cfg := configtest.NewTestConfig().Config

	cfgService := configMocks.ConfigService{}
	cfgService.On("LoadConfig").Return(cfg, nil)
	cfgService.On("NeedsMigration").Return(false, "")
	cfgService.On("UpdateAssistantConfig", mock.Anything).Return()
	cfgService.On("ConfigFilePath").Return("/foo/path")

	app := NewApp(
		WithTUI(&tui),
		WithConfigService(&cfgService),
	)

	app.EnableAssistant()

	if call := mockutil.FindMethodCall(mockutil.ShowPicker, tui.Calls); call != nil {
		assert.Equal(t, "Which AI provider for the assistant do you want to use?", call.Arguments.Get(0).(string))
		items := call.Arguments.Get(1).([]picker.Item)
		assert.Len(t, items, 5) // OpenAI, Anthropic, Gemini, Ollama, OpenAI-Compatible
		assert.Equal(t, "OpenAI", items[0].Title())
		assert.Equal(t, "Anthropic", items[1].Title())
		assert.Equal(t, "Google Gemini", items[2].Title())
		assert.Equal(t, "Ollama", items[3].Title())
		assert.Equal(t, "OpenAI-Compatible", items[4].Title())
		assert.Equal(t, "OpenAI", call.Arguments.Get(2).(*picker.Item).Title())
	}

	tui.AssertCalled(t, mockutil.Confirmation, mock.AnythingOfType("uimsg.Confirm"))
	tui.AssertCalled(t, mockutil.Print, uimsg.AssistantUpdateConfigResult(true, "/foo/path"))
}
