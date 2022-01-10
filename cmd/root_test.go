package cmd

import (
	"bytes"
	"testing"

	"emperror.dev/errors"
	"github.com/Netflix/go-expect"
	"github.com/phuslu/log"
	"github.com/stretchr/testify/assert"
)

func Test_Root(t *testing.T) {
	runVT10XCommandTest(t, []string{}, false, func(c *expect.Console, s *setup) {
		_, err := c.ExpectString(rootCmd.Long)
		assert.NoError(t, err)
	})
}

func Test_Help(t *testing.T) {
	runVT10XCommandTest(t, []string{"--help"}, false, func(c *expect.Console, s *setup) {
		_, err := c.ExpectString(rootCmd.Long)
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
