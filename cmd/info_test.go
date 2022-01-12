package cmd

import (
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	mocks "github.com/lemoony/snipkit/mocks/provider"
)

func Test_Info(t *testing.T) {
	system := testutil.NewTestSystem()
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
