package cmd

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	mocks "github.com/lemoony/snipkit/mocks/managers"
)

func Test_Info(t *testing.T) {
	configtest.ResetSnipkitHome(t)

	system := testutil.NewTestSystem()
	cfgFilePath := configtest.NewTestConfigFilePath(t, system.Fs)

	manager := mocks.Manager{}
	manager.On("Info").Return([]model.InfoLine{{Key: "Some-Key", Value: "Some-Value", IsError: false}})

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	ts := setup{
		system:   system,
		v:        v,
		provider: testProviderForManager(&manager),
	}

	runTerminalTest(t, []string{"info"}, ts, false, func(c *termtest.Console) {
		c.ExpectString("Config path: " + cfgFilePath)
		c.ExpectString("SNIPKIT_HOME: Not set")
		c.ExpectString("Theme: default")
		c.ExpectString("Some-Key: Some-Value")
	})
}
