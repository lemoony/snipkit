package termtest

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lemoony/snipkit/internal/utils/termutil"
)

type Key string

const (
	KeyEnter = Key("\r")
	KeyTab   = Key("\t")
	KeyEsc   = Key(string(rune(27)))
	KeyStrC  = Key(string(rune(3)))
	KeyUp    = Key("\u001B[A")
	KeyDown  = Key("\u001B[B")
	KeyLeft  = Key("\x1b[D")
	KeyRight = Key("\x1b[C")

	defaultTestSleepTime = time.Millisecond * 100
	defaultSendSleepTime = time.Millisecond * 10
	defaultTestTimeout   = time.Second * 2
)

func (k Key) Str() string {
	return string(k)
}

type Console struct {
	t *testing.T
	c *expect.Console
}

func (c *Console) ExpectString(val string) {
	_, err := c.c.ExpectString(val)
	assert.NoError(c.t, err)
}

func (c *Console) SendKey(val Key) {
	c.Send(string(val))
}

func (c *Console) Send(val string) {
	_, err := c.c.Send(val)
	assert.NoError(c.t, err)
	time.Sleep(defaultSendSleepTime)
}

func (c *Console) ExpectEOF() {
	_, err := c.c.ExpectEOF()
	assert.NoError(c.t, err)
}

func Keys(keys ...Key) []string {
	result := make([]string, len(keys))
	for i := range keys {
		result[i] = keys[i].Str()
	}
	return result
}

// RunTerminalTest runs a fake terminal test which catpures all in & output
// Source: https://github.com/AlecAivazis/survey/blob/master/survey_posix_test.go
func RunTerminalTest(t *testing.T, test func(c *Console), setupFunc func(termutil.Stdio)) {
	t.Helper()

	// Multiplex output to a buffer as well for the raw bytes.
	buf := new(bytes.Buffer)
	c, state, err := vt10x.NewVT10XConsole(
		expect.WithStdout(buf),
		expect.WithDefaultTimeout(defaultTestTimeout),
	)
	require.Nil(t, err)
	defer func() {
		_ = c.Close()
	}()

	procedureDone := make(chan struct{})
	go func() {
		defer close(procedureDone)
		time.Sleep(defaultTestSleepTime)
		test(&Console{t: t, c: c})
	}()

	runWithTimeout(t, func() {
		setupFunc(termutil.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()})
	}, defaultTestTimeout)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	assert.NoError(t, c.Tty().Close())
	<-procedureDone

	t.Logf("Raw output: %q", buf.String())

	// Dump the terminal's screen.
	t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))
}

func runWithTimeout(t *testing.T, procedure func(), timeout time.Duration) {
	t.Helper()

	done := make(chan struct{}, 1)
	go func() {
		defer close(done)
		procedure()
	}()

	select {
	case <-done:
		break
	case <-time.After(timeout):
		assert.Fail(t, fmt.Sprintf("function did not finish in %f seconds", timeout.Seconds()))
	}
}
