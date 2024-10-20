package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/ui/picker"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	assistantMocks "github.com/lemoony/snipkit/mocks/assistant"
	configMocks "github.com/lemoony/snipkit/mocks/config"
	managerMocks "github.com/lemoony/snipkit/mocks/managers"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_App_GenerateSnippetWithAssistant(t *testing.T) {
	const exampleScript = `
#!/bin/bash
# ${PARAM} Key: FOO_KEY
echo ${FOO_KEY}
`

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On(mockutil.ShowPrompt, mock.Anything).Return(true, "foo prompt")
	tui.On(mockutil.ShowSpinner, "foo prompt", mock.AnythingOfType("chan bool")).Return().Run(func(args mock.Arguments) {
		go func() { <-(args[1].(chan bool)) }()
	})
	tui.On(mockutil.OpenEditor, mock.Anything, mock.Anything).Return()
	tui.On(mockutil.ShowParameterForm, mock.Anything, mock.Anything, mock.Anything).Return([]string{"hello world", saveYes}, true)

	cfg := configtest.NewTestConfig().Config
	cfgService := configMocks.ConfigService{}
	cfgService.On("LoadConfig").Return(cfg, nil)
	cfgService.On("NeedsMigration").Return(false, "")

	assistantMock := assistantMocks.Assistant{}
	assistantMock.On("Query", mock.Anything).Return(exampleScript, "foo-script.sh")

	fsLibManager := managerMocks.Manager{}
	fsLibManager.On("Key").Return(fslibrary.Key)
	fsLibManager.On(mockutil.SaveAssistantSnippet, mock.Anything, mock.Anything).Return("/path", "foo-script.sh")

	provider := managerMocks.Provider{}
	provider.On("CreateManager", mock.Anything, mock.Anything).Return([]managers.Manager{&fsLibManager}, nil)

	app := NewApp(
		WithTUI(&tui),
		WithConfigService(&cfgService),
		WithProvider(&provider),
		WithAssistantProviderFunc(func(config assistant.Config) assistant.Assistant {
			return &assistantMock
		}),
	)

	app.GenerateSnippetWithAssistant()

	fsLibManager.AssertCalled(t, mockutil.SaveAssistantSnippet, "foo-script.sh", []byte(exampleScript))
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
