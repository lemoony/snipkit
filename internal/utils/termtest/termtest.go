package termtest

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	pseudotty "github.com/creack/pty"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/termutil"
)

type Key string

const (
	KeyEnter  = Key("\r")
	KeyTab    = Key("\t")
	KeyEsc    = Key(string(rune(27)))
	KeyStrC   = Key(string(rune(3)))
	KeyDelete = Key(string(rune(127)))
	KeyUp     = Key("\u001B[A")
	KeyDown   = Key("\u001B[B")
	KeyLeft   = Key("\x1b[D")
	KeyRight  = Key("\x1b[C")

	defaultTestSleepTime = time.Millisecond * 100
	defaultSendSleepTime = time.Millisecond * 10
	defaultTestTimeout   = time.Second * 5
)

// AnsiStringMatcher fulfills the Matcher interface to match strings against a given
// bytes.Buffer.
type AnsiStringMatcher struct {
	str string
}

func (sm *AnsiStringMatcher) Match(v interface{}) bool {
	buf, ok := v.(*bytes.Buffer)
	if !ok {
		return false
	}

	cleanedBufStr := buf.String()
	// cleanedBufStr := ansi.Strip(buf.String())
	if strings.Contains(cleanedBufStr, "Snip") {
		fmt.Print("x")
	}
	return strings.Contains(cleanedBufStr, sm.str)
}

func (sm *AnsiStringMatcher) Criteria() interface{} {
	return sm.str
}

func (k Key) Str() string {
	return string(k)
}

type Console struct {
	t *testing.T
	c *expect.Console
}

func (c *Console) ExpectString(val string) {
	// _, err := c.c.Expect(AnsiString(val))
	_, err := c.c.ExpectString(val)
	assert.NoError(c.t, err)
}

func AnsiString(strs ...string) expect.ExpectOpt {
	return func(opts *expect.ExpectOpts) error {
		for _, str := range strs {
			opts.Matchers = append(opts.Matchers, &AnsiStringMatcher{
				str: str,
			})
		}
		return nil
	}
}

func (c *Console) SendKey(val Key) {
	c.Send(string(val))
}

func (c *Console) Send(val string) {
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
func RunTerminalTest(t *testing.T, test func(c *Console), setupFunc func(termutil.Stdio)) {
	t.Helper()

	pty, tty, err := pseudotty.Open()
	if err != nil {
		t.Fatalf("failed to open pseudotty: %v", err)
	}

	term := vt10x.New(vt10x.WithWriter(tty))
	c, err := expect.NewConsole(expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))
	if err != nil {
		t.Fatalf("failed to create console: %v", err)
	}
	defer func() {
		_ = c.Close()
	}()

	procedureDone := make(chan struct{})
	go func() {
		defer close(procedureDone)
		time.Sleep(defaultTestSleepTime)
		test(&Console{t: t, c: c})
	}()

	runWithTimeout(
		t,
		func() { setupFunc(termutil.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}) },
		defaultTestTimeout,
	)

	<-procedureDone

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	_ = c.Tty().Close()
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
