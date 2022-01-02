package cmd

import (
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
)

func Test_Root(t *testing.T) {
	runCommandTest(t, []string{}, false, func(console *expect.Console) {
		_, err := console.ExpectString(rootCmd.Long)
		assert.NoError(t, err)
	})
}

func Test_Help(t *testing.T) {
	runCommandTest(t, []string{"--help"}, false, func(console *expect.Console) {
		_, err := console.ExpectString(rootCmd.Long)
		assert.NoError(t, err)
	})
}

func Test_UnknownCommand(t *testing.T) {
	runCommandTest(t, []string{"foo"}, true, func(console *expect.Console) {
		_, err := console.ExpectString("Error: unknown command \"foo\" for \"snipkit\"")
		assert.NoError(t, err)
	})
}
