package cmd

import (
	"bytes"
	"testing"

	"emperror.dev/errors"
	"github.com/Netflix/go-expect"
	"github.com/phuslu/log"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/config"
	appMocks "github.com/lemoony/snipkit/mocks/app"
	configMocks "github.com/lemoony/snipkit/mocks/config"
)

func Test_Root(t *testing.T) {
	runVT10XCommandTest(t, []string{}, false, func(c *expect.Console, s *setup) {
		_, err := c.ExpectString(rootCmd.Long)
		assert.NoError(t, err)
	})
}

func Test_Root_default_info(t *testing.T) {
	cfg := config.Config{}
	cfg.DefaultRootCommand = "info"

	cfgService := configMocks.ConfigService{}
	cfgService.On("LoadConfig").Return(cfg, nil)
	cfgService.On("ConfigFilePath").Return("/path/to/cfg-file")

	app := appMocks.App{}
	app.On("Info").Return()

	err := runMockedTest(t, []string{}, withConfigService(&cfgService), withApp(&app))

	assert.NoError(t, err)
	app.AssertNumberOfCalls(t, "Info", 1)
}

func Test_Help(t *testing.T) {
	runVT10XCommandTest(t, []string{"--help"}, false, func(c *expect.Console, s *setup) {
		_, err := c.ExpectString(rootCmd.Long)
		assert.NoError(t, err)
	})
}

func Test_Version(t *testing.T) {
	version := "0.0.0-SNAPSHOT-cd1c032"
	SetVersion(version)
	runVT10XCommandTest(t, []string{"--version"}, false, func(c *expect.Console, s *setup) {
		_, err := c.ExpectString("snipkit version " + version)
		assert.NoError(t, err)
	})
}

func Test_UnknownCommand(t *testing.T) {
	runVT10XCommandTest(t, []string{"foo"}, true, func(c *expect.Console, s *setup) {
		_, err := c.ExpectString("Error: unknown command \"foo\" for \"snipkit\"")
		assert.NoError(t, err)
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
