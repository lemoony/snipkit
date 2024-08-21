package cmd

import (
	"bytes"
	"path"
	"testing"

	"emperror.dev/errors"
	"github.com/phuslu/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	appMocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Root(t *testing.T) {
	defer resetCommand(rootCmd)

	system := testutil.NewTestSystem()
	cfgFilePath := configtest.NewTestConfigFilePath(t, system.Fs)

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	ts := _defaultSetup
	ts.system = system
	ts.v = v

	runTerminalTest(t, []string{}, ts, false, func(c *termtest.Console) {
		c.ExpectString(rootCmd.Long)
	})
}

func Test_Root_default_info(t *testing.T) {
	defer func() {
		configtest.ResetSnipkitHome(t)
		resetCommand(rootCmd)
	}()

	system := testutil.NewTestSystem()
	cfgFilePath := configtest.NewTestConfigFilePath(t, system.Fs, configtest.WithAdapter(func(c *config.Config) {
		c.DefaultRootCommand = "info"
	}))
	configtest.SetSnipkitHomeEnv(t, path.Dir(cfgFilePath))

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	app := appMocks.App{}
	app.On("Info").Return()

	runExecuteTest(t, []string{}, withViper(v), withSystem(system), withApp(&app))

	app.AssertNumberOfCalls(t, "Info", 1)
}

func Test_Help(t *testing.T) {
	defer resetCommand(rootCmd)

	runTerminalTest(t, []string{"--help"}, _defaultSetup, false, func(c *termtest.Console) {
		c.ExpectString(rootCmd.Long)
	})
}

func Test_Version(t *testing.T) {
	rootCmd.ResetFlags()

	version := "0.0.0-SNAPSHOT-cd1c032"
	SetVersion(version)

	runTerminalTest(t, []string{"--version"}, _defaultSetup, false, func(c *termtest.Console) {
		c.ExpectString("snipkit version " + version)
	})
}

func Test_UnknownCommand(t *testing.T) {
	runTerminalTest(t, []string{"foo"}, _defaultSetup, true, func(c *termtest.Console) {
		c.ExpectString("Error: unknown command \"foo\" for \"snipkit\"")
	})
}

func Test_handlePanic(t *testing.T) {
	tests := []struct {
		name     string
		errFunc  func()
		contains string
	}{
		{
			name:     "panic but no error",
			errFunc:  func() { panic("test") },
			contains: "Exited with panic: test",
		},
		{
			name:     "panic with error",
			errFunc:  func() { panic(errors.New("test error")) },
			contains: "Exited with panic error: test error",
		},
	}

	defer func() {
		configureLogging()
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			assert.NotPanics(t, func() {
				defer handlePanic()
				log.DefaultLogger.Writer = log.IOWriter{Writer: buf}
				log.DefaultLogger.Level = log.InfoLevel

				tt.errFunc()
			})

			assert.Contains(t, buf.String(), tt.contains)
		})
	}
}
