package picker

import (
	"bytes"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ShowPicker(t *testing.T) {
	runExpectTest(t, func(c *expect.Console) {
		_, err := c.ExpectString("Which snippet manager should be added to your configuration")
		assert.NoError(t, err)
		time.Sleep(time.Millisecond * 10)

		_, err = c.Send("\x1b[B")
		assert.NoError(t, err)
		time.Sleep(time.Millisecond * 10)

		_, err = c.Send("\x1b[B")
		assert.NoError(t, err)
		time.Sleep(time.Millisecond * 10)

		_, err = c.Send("\x1b[A")
		assert.NoError(t, err)
		time.Sleep(time.Millisecond * 10)

		_, err = c.Send("\r")
		assert.NoError(t, err)
		time.Sleep(time.Millisecond * 10)
	}, func(stdio terminal.Stdio) {
		index, ok := ShowPicker([]Item{
			NewItem("title1", "desc1"),
			NewItem("title2", "desc2"),
			NewItem("title3", "desc3"),
		}, tea.WithInput(stdio.In), tea.WithOutput(stdio.Out))
		assert.Equal(t, 1, index)
		assert.True(t, ok)
	})
}

// Source: https://github.com/AlecAivazis/survey/blob/master/survey_posix_test.go
func runExpectTest(t *testing.T, procedure func(*expect.Console), test func(terminal.Stdio)) {
	t.Helper()
	t.Parallel()

	// Multiplex output to a buffer as well for the raw bytes.
	buf := new(bytes.Buffer)
	c, state, err := vt10x.NewVT10XConsole(
		expect.WithStdout(buf),
		expect.WithDefaultTimeout(time.Second),
	)
	require.Nil(t, err)
	defer func() {
		_ = c.Close()
	}()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		time.Sleep(time.Second)
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
