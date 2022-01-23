package app

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	configMocks "github.com/lemoony/snipkit/mocks/config"
	managerMocks "github.com/lemoony/snipkit/mocks/managers"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_App_Info(t *testing.T) {
	tui := uiMocks.TUI{}
	tui.On("ApplyConfig", mock.Anything, mock.Anything).Return()

	cfg := configtest.NewTestConfig().Config

	cfgService := configMocks.ConfigService{}
	cfgService.On("LoadConfig").Return(cfg, nil)
	cfgService.On("ConfigFilePath").Return("/path/to/cfg-file")

	tui.On(mockutil.PrintMessage, mock.Anything)
	tui.On(mockutil.PrintError, mock.Anything)

	manager := managerMocks.Manager{}
	manager.On("Info").Return(model.ManagerInfo{
		Lines: []model.ManagerInfoLine{
			{Key: "Some-Key", Value: "Some-Value", IsError: false},
			{Key: "Some-Error", Value: "Some-Error", IsError: true},
		},
	})

	app := NewApp(
		WithTUI(&tui), WithConfigService(&cfgService), withManager(&manager),
	)

	app.Info()

	tui.AssertCalled(t, mockutil.PrintMessage, "Config file: /path/to/cfg-file")
	tui.AssertCalled(t, mockutil.PrintMessage, "Some-Key: Some-Value")
	tui.AssertCalled(t, mockutil.PrintError, "Some-Error: Some-Error")
}
