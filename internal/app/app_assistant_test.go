package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/ui/assistant/wizard"
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
	const exampleScript = `
#!/bin/bash
#
# Echo foo
#
# ${PARAM} Key: FOO_KEY
echo ${FOO_KEY}
`

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On(mockutil.ShowAssistantPrompt, []string{}).Return(true, "foo prompt")
	tui.On(mockutil.ShowAssistantWizard, mock.Anything).Return(true, wizard.Result{SelectedOption: wizard.OptionSaveExit, Filename: exampleFile})
	tui.On(mockutil.ShowSpinner, "Please wait, generating script...", mock.AnythingOfType("chan bool")).Return().Run(func(args mock.Arguments) {
		go func() { <-(args[1].(chan bool)) }()
	})
	tui.On(mockutil.OpenEditor, mock.Anything, mock.Anything).Return()
	tui.On(mockutil.ShowParameterForm, mock.Anything, mock.Anything, mock.Anything).Return([]string{"hello world"}, true)

	cfg := configtest.NewTestConfig().Config
	cfgService := configMocks.ConfigService{}
	cfgService.On("LoadConfig").Return(cfg, nil)
	cfgService.On("NeedsMigration").Return(false, "")

	assistantMock := assistantMocks.Assistant{}
	assistantMock.On("Query", mock.Anything).Return(exampleScript, exampleFile)
	assistantMock.On("Initialize").Return(true, uimsg.Printable{})

	fsLibManager := managerMocks.Manager{}
	fsLibManager.On("Key").Return(fslibrary.Key)
	fsLibManager.On(mockutil.SaveAssistantSnippet, mock.Anything, mock.Anything, mock.Anything).Return("/path", exampleFile)

	provider := managerMocks.Provider{}
	provider.On("CreateManager", mock.Anything, mock.Anything, &tui).Return([]managers.Manager{&fsLibManager}, nil)

	app := NewApp(
		WithTUI(&tui),
		WithConfigService(&cfgService),
		WithProvider(&provider),
		WithAssistantProviderFunc(func(config assistant.Config) assistant.Assistant {
			return &assistantMock
		}),
	)

	app.GenerateSnippetWithAssistant("", 0)

	fsLibManager.AssertCalled(t, mockutil.SaveAssistantSnippet, "Echo foo", exampleFile, []byte(exampleScript))
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
	tui.On(mockutil.ShowAssistantPrompt, []string{}).Return(true, prompt1)
	tui.On(mockutil.ShowAssistantPrompt, []string{prompt1}).Return(true, prompt2)
	tui.On(mockutil.ShowAssistantWizard, wizard.Config{ProposedFilename: exampleFile1}).Return(true, wizard.Result{SelectedOption: wizard.OptionTryAgain})
	tui.On(mockutil.ShowAssistantWizard, wizard.Config{ProposedFilename: exampleFile2}).Return(true, wizard.Result{SelectedOption: wizard.OptionDontSaveExit})
	tui.On(mockutil.ShowSpinner, "Please wait, generating script...", mock.AnythingOfType("chan bool")).Return().Run(func(args mock.Arguments) {
		go func() { <-(args[1].(chan bool)) }()
	})
	tui.On(mockutil.OpenEditor, mock.Anything, mock.Anything).Return()

	cfg := configtest.NewTestConfig().Config
	cfgService := configMocks.ConfigService{}
	cfgService.On("LoadConfig").Return(cfg, nil)
	cfgService.On("NeedsMigration").Return(false, "")

	assistantMock := assistantMocks.Assistant{}
	assistantMock.On(mockutil.Query, prompt1).Return(exampleScript1, exampleFile1)
	assistantMock.On(mockutil.Query, mock.Anything).Return(exampleScript2, exampleFile2)
	assistantMock.On(mockutil.ValidateConfig).Return(true, uimsg.Printable{})

	app := NewApp(
		WithTUI(&tui),
		WithConfigService(&cfgService),
		WithAssistantProviderFunc(func(config assistant.Config) assistant.Assistant {
			return &assistantMock
		}),
	)

	app.GenerateSnippetWithAssistant("", 0)
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
		assert.Len(t, call.Arguments.Get(1), 2)
		assert.Equal(t, call.Arguments.Get(1).([]picker.Item)[0].Title(), "OpenAI")
		assert.Equal(t, call.Arguments.Get(1).([]picker.Item)[1].Title(), "Gemini")
		assert.Equal(t, call.Arguments.Get(2).(*picker.Item).Title(), "OpenAI")
	}

	tui.AssertCalled(t, mockutil.Confirmation, mock.AnythingOfType("uimsg.Confirm"))
	tui.AssertCalled(t, mockutil.Print, uimsg.AssistantUpdateConfigResult(true, "/foo/path"))
}
