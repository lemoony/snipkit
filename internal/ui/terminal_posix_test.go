// +build !windows

package ui

import (
	"bytes"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runTest(t *testing.T, procedure func(*expect.Console), test func(terminal.Stdio)) {
	t.Helper()
	t.Parallel()

	// Multiplex output to a buffer as well for the raw bytes.
	buf := new(bytes.Buffer)
	c, state, err := vt10x.NewVT10XConsole(expect.WithStdout(buf))
	require.Nil(t, err)
	defer func() {
		_ = c.Close()
	}()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		procedure(c)
	}()

	test(terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()})

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	assert.NoError(t, c.Tty().Close())
	<-donec

	t.Logf("Raw output: %q", buf.String())

	// Dump the terminal's screen.
	t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))
}
