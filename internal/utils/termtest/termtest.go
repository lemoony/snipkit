package termtest

import (
	"bytes"
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
	KeyUp    = Key("\u001B[A")
	KeyDown  = Key("\u001B[B")
	KeyLeft  = Key("\x1b[D")
	KeyRight = Key("\x1b[C")

	defaultSendSleepTime = time.Millisecond * 10
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
	c.SendString(string(val))
}

func (c *Console) SendString(val string) {
	_, err := c.c.Send(val)
	assert.NoError(c.t, err)
	time.Sleep(defaultSendSleepTime)
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
func RunTerminalTest(t *testing.T, procedure func(test *Console), test func(termutil.Stdio)) {
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
		time.Sleep(time.Second)
		procedure(&Console{t: t, c: c})
	}()

	test(termutil.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()})

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	assert.NoError(t, c.Tty().Close())
	<-donec

	t.Logf("Raw output: %q", buf.String())

	// Dump the terminal's screen.
	t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))
}
