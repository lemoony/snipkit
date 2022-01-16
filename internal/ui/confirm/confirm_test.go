package confirm

import (
	"bytes"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func Test_Confirm(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected bool
		send     []string
	}{
		{name: "abort", expected: false, send: []string{"\n"}},
		{name: "tab / toggle", expected: true, send: []string{"\t", "\n"}},
		{name: "toggle twice", expected: false, send: []string{"\t", "\t", "\n"}},
		{name: "y", expected: true, send: []string{"y"}},
		{name: "n", expected: false, send: []string{"n"}},
		{name: "esc", expected: false, send: []string{string(rune(27))}},
		{name: "left", expected: true, send: []string{"\x1b[D", "\n"}},
		{name: "right", expected: false, send: []string{"\x1b[C", "\n"}},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			runExpectTest(t, func(c *expect.Console) {
				for _, r := range tt.send {
					_, err := c.Send(r)
					assert.NoError(t, err)
					time.Sleep(time.Millisecond * 10)
				}
			}, func(stdio terminal.Stdio) {
				result := Confirm("Are you sure?", "Hello", WithIn(stdio.In), WithOut(stdio.Out), WithFullscreen())
				assert.Equal(t, tt.expected, result)
			})
		})
	}
}

func Test_ConfirmFormatting(t *testing.T) {
	runExpectTest(t, func(c *expect.Console) {
		_, err := c.ExpectString("Hello world")
		assert.NoError(t, err)

		_, err = c.ExpectString("Are you sure?")
		assert.NoError(t, err)

		time.Sleep(time.Millisecond * 10)

		_, err = c.Send("y")
		assert.NoError(t, err)

		_, err = c.ExpectString("Yes")
		assert.NoError(t, err)
	}, func(stdio terminal.Stdio) {
		header := `Hello world`

		result := Confirm("Are you sure?", header,
			WithIn(stdio.In),
			WithOut(stdio.Out),
			WithSelectionColor("#ff0000"),
		)

		assert.True(t, result)
	})
}

// Source: https://github.com/AlecAivazis/survey/blob/master/survey_posix_test.go
func runExpectTest(t *testing.T, procedure func(*expect.Console), test func(terminal.Stdio)) {
	t.Helper()

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

func Test_zeroAwareMin(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{name: "a 1, b 1", a: 1, b: 1, expected: 1},
		{name: "a 2, b 1", a: 2, b: 1, expected: 1},
		{name: "a 1, b 2", a: 1, b: 2, expected: 1},
		{name: "a 0, b 2", a: 0, b: 2, expected: 2},
		{name: "a 2, b 0", a: 2, b: 0, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, zeroAwareMin(tt.a, tt.b))
		})
	}
}
