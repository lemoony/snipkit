package app

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snippet-kit/internal/config/configtest"
	"github.com/lemoony/snippet-kit/internal/model"
	configMocks "github.com/lemoony/snippet-kit/mocks/config"
	providerMocks "github.com/lemoony/snippet-kit/mocks/provider"
	uiMocks "github.com/lemoony/snippet-kit/mocks/ui"
)

func Test_App_Info(t *testing.T) {
	terminal := uiMocks.Terminal{}
	terminal.On("ApplyConfig", mock.Anything, mock.Anything).Return()

	cfg := configtest.NewTestConfig().Config

	cfgService := configMocks.ConfigService{}
	cfgService.On("LoadConfig").Return(cfg, nil)
	cfgService.On("ConfigFilePath").Return("/path/to/cfg-file")

	terminal.On("PrintMessage", mock.Anything)
	terminal.On("PrintError", mock.Anything)

	provider := providerMocks.Provider{}
	provider.On("Info").Return(model.ProviderInfo{
		Lines: []model.ProviderLine{
			{Key: "Some-Key", Value: "Some-Value", IsError: false},
			{Key: "Some-Error", Value: "Some-Error", IsError: true},
		},
	})

	app := NewApp(
		WithTerminal(&terminal), WithConfigService(&cfgService), withProviders(&provider),
	)

	app.Info()

	terminal.AssertCalled(t, "PrintMessage", "Config file: /path/to/cfg-file")
	terminal.AssertCalled(t, "PrintMessage", "Some-Key: Some-Value")
	terminal.AssertCalled(t, "PrintError", "Some-Error: Some-Error")
}
