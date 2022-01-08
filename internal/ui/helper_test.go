package ui

import (
	"bytes"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/gdamore/tcell/v2"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
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

func runScreenTest(t *testing.T, procedure func(s tcell.Screen), test func(s tcell.SimulationScreen)) {
	t.Helper()
	screen := mkTestScreen(t)

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		time.Sleep(time.Millisecond * 50)
		test(screen)
	}()

	procedure(screen)
	<-donec
}

func mkTestScreen(t *testing.T) tcell.SimulationScreen {
	t.Helper()
	s := tcell.NewSimulationScreen("")

	if s == nil {
		t.Fatalf("Failed to get simulation screen")
	}
	if e := s.Init(); e != nil {
		t.Fatalf("Failed to initialize screen: %v", e)
	}
	return s
}

func sendString(t *testing.T, value string, screen tcell.Screen) {
	t.Helper()
	for _, v := range value {
		assert.NoError(t, screen.PostEvent(tcell.NewEventKey(tcell.KeyRune, v, tcell.ModNone)))
		time.Sleep(time.Millisecond * 5) // sleep shortly to empty event queue
	}
}
