package cmd

import (
	"testing"

	mocks "github.com/lemoony/snipkit/mocks/config"
)

func Test_ConfigInit(t *testing.T) {
	configService := mocks.ConfigService{}
	configService.On("Create").Return()

	runExecuteTest(t, []string{"config", "init"}, withConfigService(&configService))

	configService.AssertNumberOfCalls(t, "Create", 1)
}

func Test_ConfigClean(t *testing.T) {
	configService := mocks.ConfigService{}
	configService.On("Clean").Return(nil)

	runExecuteTest(t, []string{"config", "clean"}, withConfigService(&configService))

	configService.AssertNumberOfCalls(t, "Clean", 1)
}

func Test_ConfigEdit(t *testing.T) {
	configService := mocks.ConfigService{}
	configService.On("Edit").Return(nil)

	runExecuteTest(t, []string{"config", "edit"}, withConfigService(&configService))

	configService.AssertNumberOfCalls(t, "Edit", 1)
}

func Test_ConfigMigrate(t *testing.T) {
	configService := mocks.ConfigService{}
	configService.On("Migrate").Return(nil)

	runExecuteTest(t, []string{"config", "migrate"}, withConfigService(&configService))

	configService.AssertNumberOfCalls(t, "Migrate", 1)
}
