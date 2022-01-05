package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mocks "github.com/lemoony/snippet-kit/mocks/config"
)

func Test_ConfigInit(t *testing.T) {
	configService := mocks.ConfigService{}

	configService.On("Create").Return(nil)

	err := runMockedTest(t, []string{"config", "init"}, withConfigService(&configService))

	assert.NoError(t, err)
	configService.AssertNumberOfCalls(t, "Create", 1)
}

func Test_ConfigClean(t *testing.T) {
	configService := mocks.ConfigService{}

	configService.On("Clean").Return(nil)

	err := runMockedTest(t, []string{"config", "clean"}, withConfigService(&configService))

	assert.NoError(t, err)
	configService.AssertNumberOfCalls(t, "Clean", 1)
}

func Test_ConfigEdit(t *testing.T) {
	configService := mocks.ConfigService{}

	configService.On("Edit").Return(nil)

	err := runMockedTest(t, []string{"config", "edit"}, withConfigService(&configService))

	assert.NoError(t, err)
	configService.AssertNumberOfCalls(t, "Edit", 1)
}
