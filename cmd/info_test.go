package cmd

import (
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/config/configtest"
	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/utils"
	"github.com/lemoony/snippet-kit/mocks"
)

func Test_Info(t *testing.T) {
	system := utils.NewTestSystem()
	cfgFilePath := configtest.NewTestConfigFilePath(t, system.Fs)

	provider := mocks.Provider{}
	provider.On("Info").Return(model.ProviderInfo{
		Lines: []model.ProviderLine{
			{Key: "Some-Key", Value: "Some-Value", IsError: false},
		},
	})

	runVT10XCommandTest(t, []string{"info"}, false, func(c *expect.Console, s *setup) {
		_, err := c.Expectf("Config file: %s", cfgFilePath)
		assert.NoError(t, err)
		_, err = c.ExpectString("Some-Key: Some-Value")
		assert.NoError(t, err)
	}, withSystem(system), withConfigFilePath(cfgFilePath), withProviders(&provider))
}
