package app

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	configMocks "github.com/lemoony/snipkit/mocks/config"
	providerMocks "github.com/lemoony/snipkit/mocks/provider"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_App_Info(t *testing.T) {
	terminal := uiMocks.Terminal{}
	terminal.On("ApplyConfig", mock.Anything, mock.Anything).Return()

	cfg := configtest.NewTestConfig().Config

	cfgService := configMocks.ConfigService{}
	cfgService.On("LoadConfig").Return(cfg, nil)
	cfgService.On("ConfigFilePath").Return("/path/to/cfg-file")

	terminal.On(mockutil.PrintMessage, mock.Anything)
	terminal.On(mockutil.PrintError, mock.Anything)

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

	terminal.AssertCalled(t, mockutil.PrintMessage, "Config file: /path/to/cfg-file")
	terminal.AssertCalled(t, mockutil.PrintMessage, "Some-Key: Some-Value")
	terminal.AssertCalled(t, mockutil.PrintError, "Some-Error: Some-Error")
}
