package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func runCommandTest(t *testing.T, args []string, hasError bool, test func(*expect.Console)) {
	t.Helper()

	// Multiplex output to a buffer as well for the raw bytes.
	buf := new(bytes.Buffer)
	c, state, err := vt10x.NewVT10XConsole(
		expect.WithStdout(buf),
		expect.WithDefaultTimeout(time.Second*2),
	)
	require.Nil(t, err)
	defer func() {
		_ = c.Close()
	}()

	donec := make(chan struct{})

	rootCmd.SetIn(c.Tty())
	rootCmd.SetOut(c.Tty())
	rootCmd.SetErr(c.Tty())

	defer rootCmd.ResetFlags()
	rootCmd.SetArgs(args)
	err = rootCmd.Execute()

	go func() {
		defer close(donec)
		test(c)
	}()

	<-donec

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	assert.NoError(t, c.Tty().Close())

	t.Logf("Raw output: %q", buf.String())

	// Dump the terminal's screen.
	t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))

	if hasError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}
