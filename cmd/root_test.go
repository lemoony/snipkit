package cmd

import (
	"testing"

	"github.com/Netflix/go-expect"
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
